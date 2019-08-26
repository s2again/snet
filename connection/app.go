package connection

import (
	"bytes"
	"log"
)

const gameChannel uint32 = 0

// 当前不校验session有效性，因此调用者自行保证其有效性。
func (c *Connection) SetSession(userID uint32, sessionID [16]byte) {
	c.UserID, c.SessionID = userID, sessionID
}

func (c *Connection) ListOnlineServers(getCommendList func(CommendSvrInfo)) error {
	var id MsgListenerID
	id = c.AddListener(Command_COMMEND_ONLINE, func(body bytes.Buffer) {
		c.RemoveListener(Command_COMMEND_ONLINE, id)
		info, err := parseCommendSvrInfo(&body)
		if err != nil {
			_ = c.Close()
		}
		getCommendList(info)
	})
	err := c.Send(Command_COMMEND_ONLINE, c.SessionID, gameChannel)
	if err != nil {
		return err
	}
	return nil
}

type OnlineServerInfo struct {
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
	svrList     []OnlineServerInfo
	// friendList []byte
}

func parseCommendSvrInfo(buffer *bytes.Buffer) (info CommendSvrInfo, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	log.Println("Command_COMMEND_ONLINE response bytes", buffer.Bytes())
	mustBinaryRead(buffer, &info.maxOnlineID)
	mustBinaryRead(buffer, &info.isVIP)
	mustBinaryRead(buffer, &info.onlineCnt)
	log.Println("onlineCnt", info.onlineCnt)
	info.svrList = make([]OnlineServerInfo, info.onlineCnt)
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
