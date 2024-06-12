package model

import (
	"net"
)

type ClientConn struct{}

type ConnInfo struct {
	Conn     net.Conn
	UserName string
}

var ClientConnsMap map[int]ConnInfo

func init() {
	ClientConnsMap = make(map[int]ConnInfo)
}

func (cc ClientConn) Save(userID int, name string, userConn net.Conn) {
	ClientConnsMap[userID] = ConnInfo{userConn, name}
}

func (cc ClientConn) Delete(userConn net.Conn) {
	for id, connInfo := range ClientConnsMap {
		if userConn == connInfo.Conn {
			delete(ClientConnsMap, id)
		}
	}
}

func (cc ClientConn) SearchConnectionByUserName(username string) (connInfo net.Conn, err error) {
	user, err := CurrentUserDao.GetUserByUsername(username)
	if err != nil {
		return
	}

	connInfo = ClientConnsMap[user.ID].Conn
	return
}
