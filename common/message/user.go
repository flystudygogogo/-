package message

//定义一个用户的结构体

type User struct {
	//确认字段信息
	//为了序列化和反序列化成功，必须保证用户信息的字符串和结构体的字段对应的tag，否则就会失败

	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"` //用户状态
	Sex        string `json:"sex"`        //用户性别
}
