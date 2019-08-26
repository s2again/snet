// 客户端模拟Demo
// http://51seer.61.com/?sid=
package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"main/config"
	"main/connection"
)

var (
	configFile *config.ServerConfig
	loginAddr  *net.TCPAddr

	conn *connection.Connection
)

func init() {
	var err error
	configFile, err = config.GetServerConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(configFile)
	loginAddr, err = config.GetLoginServer(configFile.IpConfig.HTTP.URL)
	if err != nil {
		panic(err)
	}
	fmt.Println(loginAddr)
}

func main() {
	conn, err := connection.Connect(loginAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Login
	var sid string
	fmt.Println("Input SID:")
	n, err := fmt.Scanf("%s", &sid)
	if n < 1 {
		return
	}
	if err != nil {
		panic(err)
	}
	login(conn, sid)
	fmt.Printf("userID: %v sessionID: %v\n", conn.UserID, conn.SessionID)

	select {}
}

func login(conn *connection.Connection, sid string) {
	userID, session, err := parseSID(sid)
	if err != nil {
		panic(err)
	}
	conn.SetSession(userID, session)
	err = conn.ListOnlineServers(func(info connection.CommendSvrInfo) {
		log.Printf("CommendSvrInfo %+v\n", info)
		go func() {
			firstOnline := info.SvrList[0]
			conn, err = loginOnline(conn.UserID, conn.SessionID, firstOnline)
			if err != nil {
				panic(err)
			}
		}()
		conn.Close()
	})
	if err != nil {
		panic(err)
	}
}

func loginOnline(userID uint32, sessionID [16]byte, server connection.OnlineServerInfo) (conn *connection.Connection, err error) {
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
	err = conn.LoginOnline(func(info connection.UserInfo) {
		fmt.Printf("UserInfo For Login %+v \n", info)
	})
	return
}

func parseSID(sid string) (userID uint32, session [16]byte, err error) {
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
