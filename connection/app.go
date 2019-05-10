package connection

import "bytes"

const gameChannel uint32 = 0

func (c *Connection) LoginWithSession(userID uint32, sid [16]byte) error {
	c.UserID, c.SID = userID, sid
	c.AddListener(Command_COMMEND_ONLINE, func(body bytes.Buffer) {
		// TODO
	})
	err := c.Send(Command_COMMEND_ONLINE, c.SID, gameChannel)
	if err != nil {
		return err
	}
	return nil
}
