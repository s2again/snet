// 客户端模拟Demo
// http://51seer.61.com/?sid=
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"main/config"
	"main/connection"
)

var (
	configFile *config.ServerConfig
	loginAddr  *net.TCPAddr
)

func init() {
	var err error
	f, err := os.OpenFile("seer.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	log.SetOutput(f)
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

	loginConn, err := connection.Connect(loginAddr)
	if err != nil {
		panic(err)
	}
	conn, err := login(loginConn, sid)
	loginConn.Close()
	fmt.Printf("userID: %v sessionID: %v\n", loginConn.UserID, loginConn.SessionID)

	petlist := mustResolvePromise(conn.GetPetList())
	fmt.Printf("精灵列表： \n")
	for _, pet := range petlist.([]connection.PetListInfo) {
		fmt.Printf("%+v\n", pet)
	}
	for _, pet := range petlist.([]connection.PetListInfo) {
		petinfo := mustResolvePromise(conn.GetPetInfo(pet.CatchTime))
		fmt.Printf("%+v\n", petinfo)
	}
	select {}
}

func login(loginConn *connection.Connection, sid string) (conn *connection.Connection, err error) {
	userID, session, err := parseSID(sid)
	if err != nil {
		panic(err)
	}
	loginConn.SetSession(userID, session)
	v, err := loginConn.ListOnlineServers().Get()
	if err != nil {
		panic(err)
	}
	info := v.(connection.CommendSvrInfo)
	firstOnline := info.SvrList[0]
	return loginOnline(loginConn.UserID, loginConn.SessionID, firstOnline)
}
