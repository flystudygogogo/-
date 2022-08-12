package process2

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/model"
	"go_code/chatroom/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn

	//增加一个字段，表示该Conn是哪个用户的
	UserId int
}

//这里我们编写通知所有在线用户的方法
//这个id要通知其他的用户，我上线了

func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历 onlineUsers，然后一个一个的发送NotifyUserStatusMes

	for id, up := range userMgr.onlineUsers {
		//过滤掉自己
		if id == userId {
			continue
		}
		//开始通知【单独写一个方法】
		up.NotifyMeToOtherOnline(userId)
	}

}

func (this *UserProcess) NotifyMeToOtherOnline(userId int) {
	//组装我们的消息NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("序列化出错..err=", err)
	}

	//将序列化后的notifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)

	//对mes再次序列化，准备发送0
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("序列化出错..err=", err)
		return
	}
	//发送，创建一个Transfer实例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeToOtherOnline err=", err)
		return
	}

}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	//1 先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterMesType
	var registerResMes message.RegisterResMes

	//我们需要到redis数据库中完成注册
	//1.先使用model.MyUserDao到redis去验证

	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "注册发生未知错误..."
		}
	} else {
		registerResMes.Code = 200
	}
	//3 将loginResMes序列胡
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("序列化失败 err=", err)
		return
	}

	//4 将data赋值给resMes
	resMes.Data = string(data)

	//5 将resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//6 发送data 我们将其封装到writePkg函数中
	//因为使用分层的模式(mvc),先创建一个Transfer 实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return

}

func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码
	//1,先从mes中去除mes.Data ,并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	//1 先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//2 再声明一个 LoginResMes
	var loginResMes message.LoginResMes

	//我们需要到redis数据库中完成验证
	//1.先使用model.MyUserDao到redis去验证

	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
		//这里我们先测试成功，然后我们再根据返回具体错误信息
	} else {
		loginResMes.Code = 200
		//这里因为用户已经登录成功，我们就把该登录成功的用户放入到userMgr中
		//将登录成功的userId赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUser(this)
		//通知其他的在线用户，我上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户的id放入到loginResMes.UsersId
		//便利UserMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登录成功")
	}

	////如果用户的id=100，密码=123456，认为合法，否则不合法
	//if loginMes.UseId == 100 && loginMes.UserPwd == "123456" {
	//	//合法
	//	loginResMes.Code = 200
	//
	//} else {
	//	//不合法
	//	loginResMes.Code = 500 //500状态码表示该用户不存在
	//	loginResMes.Error = "用户不存在，请先进行注册"
	//}

	//3 将loginResMes序列胡
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("序列化失败 err=", err)
		return
	}

	//4 将data赋值给resMes
	resMes.Data = string(data)

	//5 将resMes进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//6 发送data 我们将其封装到writePkg函数中
	//因为使用分层的模式(mvc),先创建一个Transfer 实例，然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
