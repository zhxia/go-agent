package main

import (
	"log"
	"agent/util"
	"os"
	"net"
	"encoding/json"
	"io"
	"reflect"
	//"../util"
	_"github.com/bitly/go-simplejson"
	"flag"
	"fmt"
)

var usage = `Usage:%s [options]
Options are:
	-c cmd Command type
	-d data Command body data
	-h host agent server host and port
`
var (
	cmd string
	data string
	host string
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
	}
	flag.StringVar(&host, "h", "127.0.0.1:9999", "")
	flag.StringVar(&cmd, "c", "", "")
	flag.StringVar(&data, "d", "", "")
	flag.Parse()
	byteData := []byte(data)
	client := &Client{
		Host:host,
	}
	if cmd == "upload" {
		var cmdUpload util.UploadCmd
		if err := json.Unmarshal(byteData, &cmdUpload); err != nil {
			log.Fatalln("UploadCmd json decode error:", err)
		}
		file := cmdUpload.File
		log.Println("file:", file)
		f, err := os.Stat(file)
		if err != nil {
			log.Fatalf("file stat error:%v", err)
		}
		cmdUpload.Size = f.Size()
		client.Request(util.CMD_TYPE_UPLOAD, &cmdUpload)

	} else if cmd == "list" {
		var cmdList util.ListCmd
		if err := json.Unmarshal(byteData, &cmdList); err != nil {
			log.Fatalln("ListCmd json decode error:", err)
		}
		client.Request(util.CMD_TYPE_LIST, &cmdList)
	} else if cmd == "create" {
		var cmdCreate util.CreateCmd
		if err := json.Unmarshal(byteData, &cmdCreate); err != nil {
			log.Fatalln("CreateCmd json decode error:", err)
		}
		client.Request(util.CMD_TYPE_CREATE, &cmdCreate)
	} else if cmd == "delete" {
		var cmdDelete util.DeleteCmd
		if err := json.Unmarshal(byteData, &cmdDelete); err != nil {
			log.Fatalln("DeleteCmd json decode error:", err)
		}
		client.Request(util.CMD_TYPE_DELETE, &cmdDelete)
	} else if cmd == "execute" {
		var cmdExecute util.ExecuteCmd
		if err := json.Unmarshal(byteData, &cmdExecute); err != nil {
			log.Fatalln("ExecuteCmd json decode error:", err)
		}
		client.Request(util.CMD_TYPE_EXECUTE, &cmdExecute)

	}
}

type Client struct {
	Host string
}

func (client *Client) Request(cmdType int8, command interface{}) ([]byte, error) {
	data, err := json.Marshal(command)
	strData := string(data)
	if err != nil {
		log.Fatalln("json encode error:", err)
	}
	conn, err := net.Dial("tcp", client.Host)
	defer conn.Close()
	if err != nil {
		log.Fatalf("connet error:%v\n", err)
	}
	util.SendData(conn, data, 1, cmdType)
	log.Println("send data:" + strData)

	if cmdType == util.CMD_TYPE_UPLOAD {
		filename := reflect.ValueOf(command).Elem().FieldByName("File").String()
		//发送文件
		uploadFile(filename, conn)
	}
	//获取响应
	_, respData := util.ReceiveData(conn)
	log.Println("receive data:", string(respData))
	return respData, nil
}

func uploadFile(file string, conn net.Conn) {
	buff := make([]byte, 4096)
	fi, err := os.Open(file)
	if err != nil {
		log.Fatalln("open file failed:", err)
	}
	defer fi.Close()
	for {
		n, err := fi.Read(buff)
		if err != nil && err != io.EOF {
			log.Fatalln("fie read error:", err)
		}
		if n == 0 {
			break;
		}
		conn.Write(buff[:n])
	}
	log.Println("file upload success!")
}


