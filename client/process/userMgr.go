package process

import (
	"fmt"
	"go_code/chatroom/client/model"
	"go_code/chatroom/common/message"
)

//客户端要维护的map

var onlineUsers = make(map[int]*message.User, 10)
var CurUser model.CurUser //在用户登录成功后，完成对curUser初始化

//在客户端显示在线的客户

func outputOnlineUser() {
	//遍历一把
	fmt.Println("当前在线用户列表:")
	for id := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
}

//编写一个方法，处理返回的NotifyUserStatusMes

func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	//适当的优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		//原来没有
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status

	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}
