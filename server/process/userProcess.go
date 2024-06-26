package process

import (
	"encoding/json"
	"fmt"
	"go-chat/common"
	"go-chat/server/model"
	"go-chat/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

func register(userName, passWord, passWordConfirm string) (user model.User, err error) {
	user, err = model.CurrentUserDao.Register(userName, passWord, passWordConfirm)
	return
}

func login(userName, passWord string) (user model.User, err error) {
	// 判断用户名和密码
	user, err = model.CurrentUserDao.Login(userName, passWord)
	return
}

// 响应客户端
func (up *UserProcess) responseClient(responseMessageType string, code int, data string, err error) {
	var responseMessage common.ResponseMessage
	responseMessage.Code = code
	responseMessage.Type = responseMessageType
	responseMessage.Data = data

	responseData, err := json.Marshal(responseMessage)
	if err != nil {
		fmt.Printf("responseClient序列化错误: %v", err)
	}

	dispatcher := utils.Dispatcher{Conn: up.Conn}

	dispatcher.WriteData(responseData)
}

func (up *UserProcess) UserRegister(message string) (err error) {
	var info common.RegisterMessage
	var code int
	data := ""
	err = json.Unmarshal([]byte(message), &info)
	if err != nil {
		code = common.ServerError
	}

	_, err = register(info.UserName, info.Password, info.PasswordConfirm)
	switch err {
	case nil:
		code = common.RegisterSucceed
	case model.ERROR_PASSWORD_DOES_NOT_MATCH:
		code = 402
	case model.ERROR_USER_ALREADY_EXISTS:
		code = 403
	default:
		code = 500
	}
	up.responseClient(common.RegisterResponseMessageType, code, data, err)
	return
}

func (up *UserProcess) UserLogin(message string) (err error) {
	var info common.LoginMessage
	var code int
	var data string
	err = json.Unmarshal([]byte(message), &info)
	if err != nil {
		code = common.ServerError
	}

	user, err := login(info.UserName, info.Password)

	switch err {
	case nil:
		code = common.LoginSucceed
		// save user conn status
		clientConn := model.ClientConn{}
		clientConn.Save(user.ID, user.Name, up.Conn)

		userInfo := common.UserInfo{ID: user.ID, UserName: user.Name}
		info, _ := json.Marshal(userInfo)
		data = string(info)
	case model.ERROR_USER_DOES_NOT_EXIST:
		code = 404
	case model.ERROR_USER_PWD:
		code = 403
	default:
		code = 500
	}
	up.responseClient(common.LoginResponseMessageType, code, data, err)
	return
}
