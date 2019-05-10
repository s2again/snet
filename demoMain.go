// 客户端模拟Demo
// http://51seer.61.com/?sid=
package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strconv"

	"main/config"
	"main/connection"
)

var (
	configFile *config.ServerConfig
	loginAddr  *net.TCPAddr
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
	sid := "031F02DEC2E3463BFDE809AF3443B164AD31E6A2"
	login(conn, sid)
	fmt.Println(conn.UserID, conn.Session)
	select {}
}

func login(conn *connection.Connection, sid string) {
	userID, session, err := parseSID(sid)
	if err != nil {
		panic(err)
	}
	err = conn.LoginWithSession(userID, session)
	if err != nil {
		panic(err)
	}
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
