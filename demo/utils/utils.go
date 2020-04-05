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
			prom.Reject(err)
			return
		}

		conn, err := snet.ConnectGuideServer(loginAddr)
		if err != nil {
			prom.Reject(err)
			return
		}

		conn.SetSession(uid, sessionID)
		prom.Resolve(conn)
	}()
	return prom
}

// noinspection GoUnusedFunction
func Guide2Online(conn *snet.GuideServerConnection, server snet.OnlineServerInfo) *promise.Promise {
	prom := promise.NewPromise()
	go func() {
		// login online
		onlineConn, err := snet.ConnectOnlineServer(server, conn.UserID, conn.SessionID)
		if err != nil {
			prom.Reject(err)
			return
		}
		conn.Close()
		prom.Resolve(onlineConn)
	}()
	return prom
}

func GetOnlineServerList(conn *snet.GuideServerConnection) ([]snet.OnlineServerInfo, error) {
	v, err := conn.ListOnlineServers().Get()
	if err != nil {
		return nil, err
	}
	info := v.(snet.CommendSvrInfo)
	fmt.Printf("CommendSvrInfo %+v\n", info)
	return info.SvrList, nil
}

func LoginOnlineServer(userID uint32, sessionID [16]byte, server snet.OnlineServerInfo) (conn *snet.OnlineServerConnection, err error) {
	fmt.Println("Login into Online", server.OnlineID, server.IP, server.Port)
	return snet.ConnectOnlineServer(server, userID, sessionID)
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

func AcceptAndCompleteTask(conn *snet.OnlineServerConnection, taskID uint32, param uint32) snet.NoviceFinishInfo {
	_, err := conn.AcceptTask(taskID).Get()
	if err != nil {
		panic(err)
	}
	result := MustResolvePromise(conn.CompleteTask(taskID, param))
	fmt.Println("finish novice", result)
	return result.(snet.NoviceFinishInfo)
}
