package connection

import (
	"bytes"
	"log"
	"net"
)

const ProtocolVersion byte = '1'
const packetHeadLen = 17

type MsgListener func(body bytes.Buffer)
type MsgListenerID *MsgListener

type Connection struct {
	UserID  uint32
	Session [16]byte

	tcpConn   *net.TCPConn
	listeners map[Command][]*MsgListener
	sequence  int32
}

func Connect(addr *net.TCPAddr) (conn *Connection, err error) {
	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return
	}
	conn = &Connection{
		tcpConn:  tcpConn,
		sequence: 0,
	}
	go func() {
		for {
			packet, err := depackFromStream(conn.tcpConn)
			if err != nil {
				return
			}
			if packet != nil {
				for _, listenFunc := range conn.listeners[packet.head.command] {
					if listenFunc != nil {
						(*listenFunc)(packet.body)
					}
				}
			}
		}
	}()
	return
}

func (c *Connection) Close() error {
	return c.tcpConn.Close()
}

func (c *Connection) AddListener(cmd Command, listen MsgListener) MsgListenerID {
	if c.listeners == nil {
		c.listeners = make(map[Command][]*MsgListener)
	}
	c.listeners[cmd] = append(c.listeners[cmd], &listen)
	return &listen
}

func (c *Connection) RemoveListener(cmd Command, listenID MsgListenerID) {
	if c.listeners == nil {
		return
	}
	for i, p := range c.listeners[cmd] {
		if p == listenID {
			log.Println("remove ID", listenID)
			c.listeners[cmd] = append(c.listeners[cmd][:i], c.listeners[cmd][i+1:]...)
			return
		}
	}
}

// data must be fixed-size type
func (c *Connection) Send(cmd Command, body ...interface{}) error {
	bodyBin, err := var2binary(body...)
	if err != nil {
		return err
	}
	if cmd > 1000 {
		c.sequence++
	}
	head := packetHead{
		length:   packetHeadLen + uint32(bodyBin.Len()),
		version:  ProtocolVersion,
		command:  cmd,
		userID:   c.UserID,
		sequence: c.sequence,
	}
	headBin, err := head2binary(head)
	if err != nil {
		return err
	}

	packetBin := bytes.NewBuffer(headBin.Bytes())
	_, err = packetBin.ReadFrom(bodyBin)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(packetBin.Bytes())
	if err != nil {
		return err
	}
	return nil
}
