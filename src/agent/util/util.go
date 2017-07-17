package util

import (
	"net"
	"bytes"
	"encoding/binary"
	"log"
)

const HEADER_LENGTH = 6

const CMD_TYPE_UPLOAD = 1

const CMD_TYPE_EXECUTE = 2;

const CMD_TYPE_CREATE = 3;

const CMD_TYPE_DELETE = 4;

const CMD_TYPE_LIST = 5;

type ReqHeader struct {
	Length  int32
	Version int8
	CmdType int8
}

type ResResult struct {
	Code int16        `json:"code"`
	Msg  string        `json:"msg"`
	Data string        `json:"data"`
}

type ExecuteCmd struct {
	Shell string        `json:"shell"`
	Args  []string        `json:"args"`
}

type CreateCmd struct {
	Path  string        `json:"path"`
	IsDir bool        `json:"is_dir"`
}

type UploadCmd struct {
	Size     int64 `json:"size"`
	Override bool        `json:"override"`
	File     string        `json:"file"`
}

type DeleteCmd struct {
	Files []string `json:"files"`
}

type ListCmd struct {
	Path string `json:"path"`
}

func SendData(conn net.Conn, data []byte, ver int8, cmdType int8) {
	strData := string(data)
	header := &ReqHeader{
		Length:int32(len(strData)),
		Version:ver,
		CmdType:cmdType,
	}
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, header)
	finalData := string(buffer.Bytes()) + strData
	conn.Write([]byte(finalData))
}

func ReceiveData(conn net.Conn) (ReqHeader, []byte) {
	//读取头信息
	headerBuffer := make([]byte, HEADER_LENGTH)
	n, _ := conn.Read(headerBuffer);
	if int32(n) != HEADER_LENGTH {
		log.Fatalln("header length error")
	}
	var reqHeader ReqHeader
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, headerBuffer)
	binary.Read(buffer, binary.BigEndian, &reqHeader)
	bodyBuffer := make([]byte, reqHeader.Length)
	//读取主体信息
	m, _ := conn.Read(bodyBuffer)
	if int32(m) != reqHeader.Length {
		log.Fatalln("content lenght error")
	}
	return reqHeader, bodyBuffer
}