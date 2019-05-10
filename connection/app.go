package connection

import (
	"bytes"
	"log"
)

const gameChannel uint32 = 0

func (c *Connection) LoginWithSession(userID uint32, sid [16]byte) error {
	c.UserID, c.Session = userID, sid
	var id MsgListenerID
	id = c.AddListener(Command_COMMEND_ONLINE, func(body bytes.Buffer) {
		c.RemoveListener(Command_COMMEND_ONLINE, id)
		log.Println(body.Bytes())
	})
	err := c.Send(Command_COMMEND_ONLINE, c.Session, gameChannel)
	if err != nil {
		return err
	}
	return nil
}
