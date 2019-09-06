// Implements classicalSeer
package connection

import (
	"net"

	"main/connection/core"
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
