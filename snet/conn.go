// implements classical seer network protocol
package snet

import (
	"net"

	"main/snet/core"
)

const gameChannel uint32 = 0

type Connection struct {
	*core.Connection
}

func Connect(addr *net.TCPAddr) (conn *Connection, err error) {
	coreConn, err := core.Connect(addr)
	if err != nil {
		return nil, err
	}
	return &Connection{coreConn}, nil
}

func (c *Connection) FinishTask(taskID int32) error {
	err := c.Send(Command_COMPLETE_TASK, taskID, 1)
	if err != nil {
		return err
	}
	return nil
}
