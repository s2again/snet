package snet

import (
	"log"

	"main/snet/core"
)

// echo测试。服务器原样回复客户端的body
func (c *Connection) Test(callback func(core.PacketBody), data ...interface{}) error {
	v, err := c.SendInPromise(Command_TEST, data...).Get()
	if err != nil {
		log.Println("Test promise rejected: ", err)
		return err
	}
	callback(v.(core.PacketBody))
	return nil
}
