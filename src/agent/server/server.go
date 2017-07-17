package main

import (
	"log"
	"net"
	"agent/util"
	"encoding/json"
	"os"
	"path"
	//"../util"
	"path/filepath"
	"strings"
	"os/exec"
	"flag"
	"fmt"
)

var usage = `Usage:%s [options]
	Options are:
		-d dataroot Set data root
		-s shellroot Set shell script root
		-h hostport Set agent listen host and port
`
var (
	dataroot string
	shellroot string
	host string
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.StringVar(&dataroot, "d", "/data", "")
	flag.StringVar(&shellroot, "s", "/data/script", "")
	flag.StringVar(&host, "h", "127.0.0.1:9999", "")
	flag.Parse()
	server := Server{
		Host:host,
	}
	server.Run()
}

type Server struct {
	Host string
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.Host)
	defer listener.Close()
	if err != nil {
		log.Fatalf("listen error:%v\n", err.Error())
	}
	log.Println("server is running...")
	log.Println("dataroot:", dataroot, ",shellroot:", shellroot)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error:%v\n", err.Error())
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reqHeader, data := util.ReceiveData(conn)
	log.Println("receive data,command type:", reqHeader.CmdType, ",data:", string(data))
	process(reqHeader.CmdType, data, conn)
	content := []byte("request done")
	util.SendData(conn, content, reqHeader.Version, 0)
}

func process(cmdType int8, data []byte, conn net.Conn) {
	if cmdType == util.CMD_TYPE_UPLOAD {
		//上传文件
		var cmdUpload util.UploadCmd
		err := json.Unmarshal(data, &cmdUpload)
		if err != nil {
			log.Println("UploadCmd json decode error!")
			return
		}
		log.Printf("file:%s", cmdUpload.File)
		log.Printf("file size:%d", cmdUpload.Size)
		buffer := make([]byte, 4096)
		total := int64(0);
		filename := path.Base(cmdUpload.File)
		destFile := dataroot + "/" + filename
		fout, err := os.Create(destFile)
		defer fout.Close()
		if err != nil {
			log.Println("create file:%s failed!", destFile)
			return
		}
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				log.Println("receive file error:", err)
				return
			}
			fout.Write(buffer[:n])
			total += int64(n)
			if total >= cmdUpload.Size {
				break;
			}
		}
	} else if cmdType == util.CMD_TYPE_CREATE {
		//创建目录
		var createCmd util.CreateCmd
		err := json.Unmarshal(data, &createCmd)
		if err != nil {
			log.Println("CreateCmd json decode error!")
			return
		}
		if createCmd.IsDir {
			os.MkdirAll(createCmd.Path, 0755)
		} else {
			dir := path.Dir(createCmd.Path)
			if err = os.MkdirAll(dir, 0755); err != nil {
				log.Println("dir created failed!")
				return
			}
			os.Create(createCmd.Path)
		}

	} else if cmdType == util.CMD_TYPE_DELETE {
		//删除文件
		var deleteCmd util.DeleteCmd
		if err := json.Unmarshal(data, &deleteCmd); err != nil {
			log.Println("DeleteCmd json decode error!")
			return
		}
		for _, f := range deleteCmd.Files {
			err := os.Remove(f)
			if err != nil {
				util.SendData(conn, []byte(err.Error()), 1, 0)
			}
		}
	} else if cmdType == util.CMD_TYPE_EXECUTE {
		var executeCmd util.ExecuteCmd
		if err := json.Unmarshal(data, &executeCmd); err != nil {
			log.Println("ExecuteCmd json decode error!")
			return
		}
		if strings.HasPrefix(executeCmd.Shell, shellroot) {
			cmd := exec.Command(executeCmd.Shell, executeCmd.Args...)
			output, err := cmd.Output()
			if err != nil {
				log.Println("get shell script execute result error:", err)
				return
			}
			log.Println("output:", string(output))
			util.SendData(conn, output, 1, 0)
		} else {
			log.Println("invalid shell script file:", executeCmd.Shell)
		}

	} else if cmdType == util.CMD_TYPE_LIST {
		var listCmd util.ListCmd
		err := json.Unmarshal(data, &listCmd)
		if err != nil {
			log.Println("ListCmd json decode error")
		}
		files := listDir(listCmd.Path)
		data, _ := json.Marshal(files)
		util.SendData(conn, data, 1, 0)
	}
}

func listDir(pathRoot string) []string {
	var files []string
	err := filepath.Walk(pathRoot, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Println("ListCmd error:", err)
	}
	return files
}
