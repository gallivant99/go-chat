package common

const (
	LoginMessageType             = "LoginMessage"
	RegisterMessageType          = "RegisterMessage"
	RegisterResponseMessageType  = "ResponseMessageType"
	LoginResponseMessageType     = "LoginResponseMessageType"

	ServerError = 500

	RegisterSucceed  = 200

	LoginSucceed = 200
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type ResponseMessage struct {
	Type  string
	Code  int    // 404 用户没找到， 403 账号或者密码错误, 200 登陆成功, 500 服务端错误
	Error string // 错误消息
	Data  string
}

type LoginMessage struct {
	UserName string
	Password string
}

type RegisterMessage struct {
	UserName        string
	Password        string
	PasswordConfirm string
}

// on line user info
type UserInfo struct {
	ID       int
	UserName string
}
