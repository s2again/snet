// implements classical seer network protocol
package snet

import (
	"net"
	"strconv"

	"github.com/s2again/snet/core"
)

const gameChannel uint32 = 0

type Connection struct {
	*core.Connection
}

type GuideServerConnection struct {
	Connection
}
type OnlineServerConnection struct {
	Connection
}

func Connect(addr *net.TCPAddr) (conn *Connection, err error) {
	coreConn, err := core.Connect(addr)
	if err != nil {
		return nil, err
	}
	return &Connection{coreConn}, nil
}

func ConnectGuideServer(addr *net.TCPAddr) (conn *GuideServerConnection, err error) {
	c, e := Connect(addr)
	if e != nil {
		return nil, e
	}
	return &GuideServerConnection{*c}, nil
}

func ConnectOnlineServer(server OnlineServerInfo, userID uint32, sessionID [16]byte) (conn *OnlineServerConnection, err error) {
	addrStr := server.IP + ":" + strconv.Itoa(int(server.Port))
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}
	c, err := Connect(addr)
	if err != nil {
		return nil, err
	}
	c.SetSession(userID, sessionID)
	return &OnlineServerConnection{*c}, nil
}
