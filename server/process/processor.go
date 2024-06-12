package process

import (
	"fmt"
	"go-chat/common"
	"go-chat/server/model"
	"go-chat/server/utils"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 处理消息
// 根据消息的类型，使用对应的处理方式
func (p *Processor) messageProcess(message common.Message) (err error) {
	switch message.Type {
	case common.LoginMessageType:
		up := UserProcess{Conn: p.Conn}
		err = up.UserLogin(message.Data)
		if err != nil {
			fmt.Printf("some error: %v\n", err)
		}
	case common.RegisterMessageType:
		up := UserProcess{Conn: p.Conn}
		err = up.UserRegister(message.Data)
		if err != nil {
			fmt.Printf("some error when register: %v\n", err)
		}
	default:
		fmt.Printf("other type\n")
	}
	return
}

// 处理和用户的之间的通讯
func (p *Processor) MainProcess() {

	// 循环读来自客户端的消息
	for {
		dispatcher := utils.Dispatcher{Conn: p.Conn}
		message, err := dispatcher.ReadData()
		if err != nil {
			if err == io.EOF {
				cc := model.ClientConn{}
				cc.Delete(p.Conn)
				fmt.Printf("client closed!\n")
				break
			}
			fmt.Printf("get login message error: %v", err)
		}

		// 处理来客户端的消息
		// 按照消息的类型，使用不同的处理方法
		err = p.messageProcess(message)
		if err != nil {
			fmt.Printf("some error: %v\n", err)
			break
		}
	}
}
