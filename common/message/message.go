package message

//确定一些消息类型

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatus"
	SmsMesType              = "SmsMes"
)

//这里我们定义几个用户状态的常量

const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息的内容
}

//先定义两个消息，后续需要在添加

type LoginMes struct {
	UserId   int    `json:"useId"`    //用户ID
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code    int    `json:"code"` //状态码 500表示该用户未注册 200表示登录成功
	UsersId []int  //增加字段，保存用户的切片
	Error   string `json:"error"` //返回的错误信息，如果没有则为nil
}

type RegisterMes struct {
	User User `json:"user"` //类型就是User的结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //状态码 400表示该用户已经占用 200表示注册成功
	Error string `json:"error"` //返回的错误信息，如果没有则为nil
}

//为了配合服务器端推送用户状态变化的消息

type NotifyUserStatusMes struct {
	UserId int `json:"userId"` //用户id
	Status int `json:"status"` //用户的状态

}

//增加一个SmsMes //发送的消息

type SmsMes struct {
	Content string `json:"content"` //内容
	User           //匿名的结构体
}
