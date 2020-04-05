// unavailable
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"main/config"
	"main/snet"
)

var (
	configFile *config.ServerConfig
	guideAddr  *net.TCPAddr
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
	guideAddr, err = configFile.GetGuideServerByHTTP()
	if err != nil {
		panic(err)
	}
	fmt.Println(guideAddr)
}

func main() {
	loginConn, err := snet.Connect(guideAddr)
	if err != nil {
		panic(err)
	}
	loginConn.SendInPromise(103)
}
