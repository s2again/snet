package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/fanliao/go-promise"

	"main/connection"
)

func LoginOnline(userID uint32, sessionID [16]byte, server connection.OnlineServerInfo) (conn *connection.Connection, err error) {
	addrStr := server.IP + ":" + strconv.Itoa(int(server.Port))
	fmt.Println("Login into Online", addrStr)
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return
	}
	conn, err = connection.Connect(addr)
	if err != nil {
		return
	}
	conn.SetSession(userID, sessionID)
	err = conn.LoginOnlineAndCallback(func(info connection.ResponseForLogin) {
		fmt.Printf("ResponseForLogin For Login %+v \n", info)
	})
	return
}

func ParseSID(sid string) (userID uint32, session [16]byte, err error) {
	if len(sid) != 40 {
		err = errors.New("illegal parameter")
		return
	}
	userIDtmp, err := strconv.ParseUint(sid[:8], 16, 32)
	userID = uint32(userIDtmp)

	sessiontmp, err := hex.DecodeString(sid[8:40])
	if err != nil {
		return
	}
	copy(session[:], sessiontmp[:32])
	return
}

func MustResolvePromise(p *promise.Promise) interface{} {
	v, err := p.Get()
	if err != nil {
		panic(err)
	}
	return v
}

func AcceptAndCompleteTask(conn *connection.Connection, taskID uint32, param uint32) connection.NoviceFinishInfo {
	_, err := conn.AcceptTask(taskID).Get()
	if err != nil {
		panic(err)
	}
	result := MustResolvePromise(conn.CompleteTask(taskID, param))
	fmt.Println("finish novice", result)
	return result.(connection.NoviceFinishInfo)
}
