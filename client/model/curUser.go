package model

import (
	"go_code/chatroom/common/message"
	"net"
)

//因为在客户端，很多地方会用到，我们将其作为一个全局
type CurUser struct {
	Conn net.Conn

	message.User
}
