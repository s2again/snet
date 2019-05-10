package connection

import (
	"bytes"
	"log"
)

const gameChannel uint32 = 0

func (c *Connection) LoginWithSession(userID uint32, sid [16]byte, getCommendList func(CommendSvrInfo)) error {
	c.UserID, c.Session = userID, sid
	var id MsgListenerID
	id = c.AddListener(Command_COMMEND_ONLINE, func(body bytes.Buffer) {
		c.RemoveListener(Command_COMMEND_ONLINE, id)
		log.Println(body.Bytes())
		info, err := parseCommendSvrInfo(&body)
		if err != nil {
			_ = c.Close()
		}
		getCommendList(info)
	})
	err := c.Send(Command_COMMEND_ONLINE, c.Session, gameChannel)
	if err != nil {
		return err
	}
	return nil
}

type ServerInfo struct {
	onlineID uint32
	userCnt  uint32
	ip       string
	port     uint16
	friends  uint32
}
type CommendSvrInfo struct {
	maxOnlineID uint32
	isVIP       uint32
	onlineCnt   uint32
	svrList     []ServerInfo
	// friendList []byte
}

func parseCommendSvrInfo(buffer *bytes.Buffer) (info CommendSvrInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	mustBinaryRead(buffer, &info.maxOnlineID)
	mustBinaryRead(buffer, &info.isVIP)
	mustBinaryRead(buffer, &info.onlineCnt)
	log.Println("onlineCnt", info.onlineCnt)
	info.svrList = make([]ServerInfo, info.onlineCnt)
	for i := uint32(0); i < info.onlineCnt; i++ {
		mustBinaryRead(buffer, &info.svrList[i].onlineID)
		mustBinaryRead(buffer, &info.svrList[i].userCnt)
		{
			var ipBin [16]byte
			mustBinaryRead(buffer, &ipBin)
			info.svrList[i].ip = string(ipBin[:])
		}
		mustBinaryRead(buffer, &info.svrList[i].port)
		mustBinaryRead(buffer, &info.svrList[i].friends)
	}
	return
}
