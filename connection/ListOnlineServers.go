package connection

import (
	"bytes"
	"log"

	"main/connection/core"
)

type OnlineServerInfo struct {
	OnlineID uint32
	UserCnt  uint32
	IP       string
	Port     uint16
	Friends  uint32
}
type CommendSvrInfo struct {
	MaxOnlineID uint32
	IsVIP       uint32
	OnlineCnt   uint32
	SvrList     []OnlineServerInfo
	// friendList []byte
}

func (c *Connection) ListOnlineServers(callback func(CommendSvrInfo)) error {
	prom, err := c.SendForPromise(Command_COMMEND_ONLINE, c.SessionID, gameChannel)
	if err != nil {
		return err
	}
	prom.OnSuccess(func(v interface{}) {
		list, err := parseCommendSvrInfo(v.(core.PacketBody))
		if err != nil {
			c.Close()
			log.Println("parseCommendSvrInfo error: ", v, "connection terminated.")
		}
		callback(list)
	}).OnFailure(func(v interface{}) {
		log.Println("ListOnlineServers promise rejected: ", v)
	})
	return nil
}

func parseCommendSvrInfo(buffer core.PacketBody) (info CommendSvrInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	log.Println("Command_COMMEND_ONLINE response bytes", buffer.Bytes())
	core.MustBinaryRead(buffer, &info.MaxOnlineID)
	core.MustBinaryRead(buffer, &info.IsVIP)
	core.MustBinaryRead(buffer, &info.OnlineCnt)
	log.Println("OnlineCnt", info.OnlineCnt)
	info.SvrList = make([]OnlineServerInfo, info.OnlineCnt)
	for i := uint32(0); i < info.OnlineCnt; i++ {
		core.MustBinaryRead(buffer, &info.SvrList[i].OnlineID)
		core.MustBinaryRead(buffer, &info.SvrList[i].UserCnt)
		{
			var ipBin [16]byte
			core.MustBinaryRead(buffer, &ipBin)
			info.SvrList[i].IP = string(bytes.Trim(ipBin[:], "\u0000"))
		}
		core.MustBinaryRead(buffer, &info.SvrList[i].Port)
		core.MustBinaryRead(buffer, &info.SvrList[i].Friends)
	}
	return
}
