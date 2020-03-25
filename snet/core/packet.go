package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
)

type Command = uint32

type SendPacketHead struct {
	length   uint32
	version  byte
	command  Command
	userID   uint32
	sequence int32
}
type RecvPacketHead struct {
	length  uint32
	version byte
	command Command
	userID  uint32
	errno   int32
}

// *bytes.Buffer弱化版接口
type PacketBody interface {
	Bytes() []byte
	Len() int
	Truncate(n int)
	Read(p []byte) (n int, err error)
	Next(n int) []byte
	ReadByte() (byte, error)
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
	UnreadByte() error
	ReadBytes(delim byte) (line []byte, err error)
	ReadString(delim byte) (line string, err error)
}
type SendPacket struct {
	head SendPacketHead
	body PacketBody
}
type RecvPacket struct {
	head RecvPacketHead
	body PacketBody
}

func depackFromStream(reader io.Reader) (pack *RecvPacket, err error) {
	const maxPacketLength = 65536
	var buffer [maxPacketLength]byte

	index := 0

	for index != packetHeadLen {
		n, err := reader.Read(buffer[index:packetHeadLen]) // receive head
		if err != nil {
			log.Println(err)
			return nil, err
		}
		log.Println("response bytes", buffer[:n])
		index += n
		log.Printf("Receive Packet Pead %d/%d bytes\n", index, packetHeadLen)
	}
	log.Println("Packet Head: ", buffer[:packetHeadLen])
	head, err := parseHead(bytes.NewReader(buffer[:packetHeadLen]))
	if err != nil {
		return nil, err
	}
	log.Printf("Parse Head %+v", head)
	if head.length > maxPacketLength {
		err = errors.New(fmt.Sprintf("Too Large Packet(%d bytes)", head.length))
		log.Println(err.Error())
		return nil, err
	}

	// index == packetHeadLen
	bodyLen := head.length - packetHeadLen
	for index != int(head.length) {
		n, err := reader.Read(buffer[index:head.length]) // receive body
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		index += n
		log.Printf("Receive Body %d/%d bytes \n", index-packetHeadLen, bodyLen)
	}
	log.Printf("Packet Body (total %d bytes) %X\n", bodyLen, buffer[packetHeadLen:head.length])
	var body bytes.Buffer
	body.Write(buffer[packetHeadLen:head.length])
	pack = &RecvPacket{
		head: head,
		body: &body,
	}
	return
}

func parseHead(input *bytes.Reader) (head RecvPacketHead, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	MustBinaryRead(input, &head.length)
	MustBinaryRead(input, &head.version)
	MustBinaryRead(input, &head.command)
	MustBinaryRead(input, &head.userID)
	MustBinaryRead(input, &head.errno)
	return
}
