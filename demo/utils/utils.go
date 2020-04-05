package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/fanliao/go-promise"

	"main/snet"
	"main/snet/core"
)

// noinspection GoUnusedFunction
func ConnectGuideServer(loginAddr *net.TCPAddr, sid string) (prom *promise.Promise) {
	prom = promise.NewPromise()
	go func() {
		defer func() {
			x := recover()
			if x != nil {
				err, ok := x.(error)
				if !ok {
					err = errors.New("promise rejected: " + fmt.Sprint(x))
				}
				prom.Reject(err)
			}
		}()
		uid, sessionID, err := core.ParseSIDString(sid)
		if err != nil {
			panic(err)
		}

		conn, err := snet.Connect(loginAddr)
		if err != nil {
			panic(err)
		}

		conn.SetSession(uid, sessionID)
		prom.Resolve(conn)
	}()
	return prom
}

// noinspection GoUnusedFunction
func Guide2Online(conn *snet.Connection, server snet.OnlineServerInfo) *promise.Promise {
	prom := promise.NewPromise()
	go func() {
		// login online
		addrStr := server.IP + ":" + strconv.Itoa(int(server.Port))
		fmt.Println("Login into Online", addrStr)
		addr, err := net.ResolveTCPAddr("tcp", addrStr)
		if err != nil {
			panic(err)
		}

		onlineConn, err := snet.Connect(addr)
		if err != nil {
			panic(err)
		}
		onlineConn.SetSession(conn.UserID, conn.SessionID)
		conn.Close()
		prom.Resolve(onlineConn)
	}()
	return prom
}

func GetOnlineServerList(subConn *snet.Connection) ([]snet.OnlineServerInfo, error) {
	v, err := subConn.ListOnlineServers().Get()
	if err != nil {
		return nil, err
	}
	info := v.(snet.CommendSvrInfo)
	fmt.Printf("CommendSvrInfo %+v\n", info)
	return info.SvrList, nil
}

func LoginOnlineServer(userID uint32, sessionID [16]byte, server snet.OnlineServerInfo) (conn *snet.Connection, err error) {
	addrStr := server.IP + ":" + strconv.Itoa(int(server.Port))
	fmt.Println("Login into Online", addrStr)
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return
	}
	conn, err = snet.Connect(addr)
	if err != nil {
		return
	}
	conn.SetSession(userID, sessionID)
	err = conn.LoginOnlineAndCallback(func(info snet.ResponseForLogin) {
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

func AcceptAndCompleteTask(conn *snet.Connection, taskID uint32, param uint32) snet.NoviceFinishInfo {
	_, err := conn.AcceptTask(taskID).Get()
	if err != nil {
		panic(err)
	}
	result := MustResolvePromise(conn.CompleteTask(taskID, param))
	fmt.Println("finish novice", result)
	return result.(snet.NoviceFinishInfo)
}
