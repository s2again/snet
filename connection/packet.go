package connection

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
)

type packetHead struct {
	length   uint32
	version  byte
	command  Command
	userID   uint32
	sequence int32
}

// *bytes.Buffer弱化版接口
type packetBody interface {
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
type packet struct {
	head packetHead
	body packetBody
}

func depackFromStream(reader io.Reader) (pack *packet, err error) {
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
	pack = &packet{
		head: head,
		body: &body,
	}
	return
}

func parseHead(input *bytes.Reader) (head packetHead, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	mustBinaryRead(input, &head.length)
	mustBinaryRead(input, &head.version)
	mustBinaryRead(input, &head.command)
	mustBinaryRead(input, &head.userID)
	mustBinaryRead(input, &head.sequence)
	return
}
