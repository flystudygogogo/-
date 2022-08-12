package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据...")
	//conn.read 在conn没有被关闭的情况下才会堵塞
	//如果客户端关闭了conn，则不会堵塞
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])

	//根据pkgLen 读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//fmt.Println("conn read fail err=", err)
		return
	}

	//把pkgLen反序列化成message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes) //注意&符号一定要加
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}
