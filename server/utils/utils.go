package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"go-chat/common"
	"net"
)

type Dispatcher struct {
	Conn net.Conn
	Buf  [10240]byte
}

func (dp *Dispatcher) ReadData() (message common.Message, err error) {
	// 读取消息长度信息
	n, err := dp.Conn.Read(dp.Buf[:4])
	if err != nil || n != 4 {
		return
	}
	var dataLen uint32 = binary.BigEndian.Uint32(dp.Buf[:4])

	// 读取消息本身
	n, err = dp.Conn.Read(dp.Buf[:dataLen])
	if err != nil {
		fmt.Printf("ReadData消息读取错误: %v", err)
		return
	}

	// 对比消息本身的长度和期望长度是否匹配
	if n != int(dataLen) {
		err = errors.New("消息长度不匹配！")
		return
	}

	// 从 conn 中解析消息并存放到 message 中，此处一定传递的是 message 的地址
	err = json.Unmarshal(dp.Buf[:dataLen], &message)
	if err != nil {
		fmt.Printf("ReadData反序列化错误: %v", err)
	}
	return
}

func (dp *Dispatcher) WriteData(data []byte) (err error) {
	var dataLen uint32 = uint32(len(data))
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[:4], dataLen)

	// 将消息长度发送给客户端
	_, err = dp.Conn.Write(bytes[:])
	if err != nil {
		fmt.Printf("WriteData发送数据长度错误: %v\n", err)
		return
	}

	// 发送消息本身给客户端
	_, err = dp.Conn.Write(data)
	if err != nil {
		fmt.Printf("WriteData发送数据错误: %v", err)
		return
	}
	return
}