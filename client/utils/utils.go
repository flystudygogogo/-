package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"net"
)

//这里将这些方法关联到结构体中

type Transfer struct {
	//分析应该有哪些字段
	Conn net.Conn
	Buf  [8096]byte //这时传输时，使用缓冲
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据...")
	//conn.read 在conn没有被关闭的情况下才会堵塞
	//如果客户端关闭了conn，则不会堵塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}
	//根据buf[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])

	//根据pkgLen 读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//fmt.Println("conn read fail err=", err)
		return
	}

	//把pkgLen反序列化成message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes) //注意&符号一定要加
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) err=", err)
		return
	}
	//发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) err=", err)
		return
	}
	return
}
