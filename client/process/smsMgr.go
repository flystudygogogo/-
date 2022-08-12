package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
)

func outputGroupMes(mes *message.Message) {
	//显示即可
	//1.反序列化mes.Data
	var smeMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smeMes)
	if err != nil {
		fmt.Println("err=", err.Error())
		return
	}

	//显示内容
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s", smeMes.UserId, smeMes.Content)
	fmt.Println(info)

}
