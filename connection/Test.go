package connection

import (
	"log"

	"main/connection/core"
)

// echo测试。服务器原样回复客户端的body
func (c *Connection) Test(parseFunc func(core.PacketBody), data ...interface{}) error {
	prom, err := c.SendForPromise(Command_Test, data...)
	if err != nil {
		return err
	}
	prom.OnSuccess(func(v interface{}) {
		parseFunc(v.(core.PacketBody))
	}).OnFailure(func(v interface{}) {
		log.Println("Test promise rejected: ", v)
	})
	return nil
}
