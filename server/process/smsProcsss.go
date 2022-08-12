package process2

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/utils"
	"net"
)

type SmsProcess struct {
}

//写方法转发消息

func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的map
	//将消息转发出去

	//取出mes的内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("err=", err)
		return
	}
	data, err := json.Marshal(mes)

	if err != nil {
		return
	}
	for id, up := range userMgr.onlineUsers {

		if id == smsMes.UserId {
			continue
		}
		this.SendMesToEachOnlineUser(data, up.Conn)
	}
}

func (this *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建一个Transfer实例，

	tf := &utils.Transfer{
		Conn: conn,
	}

	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
		return
	}
}
