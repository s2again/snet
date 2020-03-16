// Implements Taomee Network Protocol
package core

import (
	"bytes"
	"errors"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/fanliao/go-promise"
	"github.com/juju/utils/deque"
)

const ProtocolVersion byte = '1'
const packetHeadLen = 17

type MsgListener func(body PacketBody)
type MsgListenerID *MsgListener

type Connection struct {
	UserID    uint32
	SessionID [16]byte
	sequence  int32

	listeners            map[Command][]*MsgListener
	listenersLock        sync.RWMutex
	responsePromises     hashmap.Map
	responsePromisesLock sync.RWMutex

	tcpConn *net.TCPConn
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
	conn.responsePromises.Clear()
	go conn.handlePacket()
	return
}

// 当前不校验session有效性，因此调用者自行保证其有效性。
func (c *Connection) SetSession(userID uint32, sessionID [16]byte) {
	c.UserID, c.SessionID = userID, sessionID
}

func (c *Connection) handlePacket() {
	for {
		packet, err := depackFromStream(c.tcpConn)
		if err != nil {
			return
		}
		if packet != nil {
			// find listeners
			listeners := c.matchListeners(packet.head.command)
			for _, f := range listeners {
				(*f)(packet.body)
			}
			// find promise
			firstPromise := c.firstPromise(packet.head.command)
			if firstPromise != nil {
				if packet.head.errno != 0 {
					err := firstPromise.Reject(errors.New("Errno" + strconv.Itoa(int(packet.head.errno))))
					if err != nil {
						log.Println("Reject Promise Failed: ", err)
					}
				} else {
					err := firstPromise.Resolve(packet.body)
					if err != nil {
						log.Println("Resolve Promise Failed: ", err)
					}
				}
			}
			drop := len(listeners) == 0 && firstPromise == nil
			if drop {
				log.Printf("Unhandled Packet %+v\n", packet)
			}
		}
	}
}

// Get all the listeners that match command
func (c *Connection) matchListeners(cmd Command) []*MsgListener {
	c.listenersLock.RLock()
	defer c.listenersLock.RUnlock()
	var matchListeners []*MsgListener
	for _, listenFunc := range c.listeners[cmd] {
		if listenFunc != nil {
			// 先保存，解锁后再调用。如果直接调用，因为用户函数可能调用本库其他函数，造成重复锁而死锁
			matchListeners = append(matchListeners, listenFunc)
		}
	}
	return matchListeners
}

// Get the first promise in the queue correlate with cmd,
// and remove first promise from the queue.
func (c *Connection) firstPromise(cmd Command) *promise.Promise {
	c.responsePromisesLock.Lock()
	defer c.responsePromisesLock.Unlock()
	val, found := c.responsePromises.Get(cmd)
	if !found {
		return nil
	}
	promDeque := val.(*deque.Deque)
	val, has := promDeque.PopFront()
	if !has {
		return nil
	}
	return val.(*promise.Promise)
}

func (c *Connection) Close() error {
	return c.tcpConn.Close()
}

func (c *Connection) AddListener(cmd Command, listen MsgListener) MsgListenerID {
	c.listenersLock.Lock()
	defer c.listenersLock.Unlock()
	if c.listeners == nil {
		c.listeners = make(map[Command][]*MsgListener)
	}
	c.listeners[cmd] = append(c.listeners[cmd], &listen)
	id := &listen
	log.Println("add Listener", cmd, id)
	return id
}

func (c *Connection) RemoveListener(cmd Command, listenID MsgListenerID) {
	if c.listeners == nil {
		return
	}
	c.listenersLock.Lock()
	defer c.listenersLock.Unlock()
	for i, p := range c.listeners[cmd] {
		if p == listenID {
			log.Println("remove Listener", listenID)
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
	head := SendPacketHead{
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
	log.Printf("Send Message %+v %v \n", head, packetBin.Bytes())
	_, err = c.tcpConn.Write(packetBin.Bytes())
	if err != nil {
		return err
	}
	return nil
}

// DEPRECATED:
// data must be fixed-size type
// err only reports errors in sending operation. Use promise.OnFailure if you want to check response errors
func (c *Connection) SendForPromise(cmd Command, body ...interface{}) (responsePromise *promise.Promise, err error) {
	c.responsePromisesLock.Lock()
	defer c.responsePromisesLock.Unlock()
	err = c.Send(cmd, body...)
	if err != nil {
		return
	}
	val, found := c.responsePromises.Get(cmd)
	if !found {
		val = deque.New()
		c.responsePromises.Put(cmd, val)
	}
	responsePromise = promise.NewPromise()
	promDeque := val.(*deque.Deque)
	promDeque.PushBack(responsePromise)
	return
}

// data must be fixed-size type
// all errors will be reported via promise. Use promise.OnFailure if you want to check errors
//
func (c *Connection) SendInPromise(cmd Command, body ...interface{}) (responsePromise *promise.Promise) {
	responsePromise = promise.NewPromise()
	c.responsePromisesLock.Lock()
	defer c.responsePromisesLock.Unlock()
	log.Printf("send body: %X\n", body)
	err := c.Send(cmd, body...)
	if err != nil {
		responsePromise.Reject(err)
		return
	}
	promDeque, found := c.responsePromises.Get(cmd)
	if !found {
		promDeque = deque.New()
		c.responsePromises.Put(cmd, promDeque)
	}
	promDeque.(*deque.Deque).PushBack(responsePromise)
	return
}
