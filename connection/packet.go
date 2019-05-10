package connection

import (
	"bytes"
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
type packetBody = bytes.Buffer
type packet struct {
	head packetHead
	body packetBody
}

func depackFromStream(reader io.Reader) (pack *packet, err error) {
	const maxPacketLength = 65536
	var buffer [maxPacketLength]byte

	n, err := reader.Read(buffer[:packetHeadLen]) // receive head
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if n != packetHeadLen {
		log.Println("Only Receive Packet Head Bytes Length", n)
		return nil, err
	}
	head, err := parseHead(bytes.NewReader(buffer[:packetHeadLen]))
	if err != nil {
		return nil, err
	}
	n, err = reader.Read(buffer[packetHeadLen : head.length-packetHeadLen]) // receive body
	if err != nil {
		return nil, err
	}
	fmt.Printf("%X", buffer[:head.length])
	var body bytes.Buffer
	body.Write(buffer[packetHeadLen : head.length-packetHeadLen])
	pack = &packet{
		head: head,
		body: packetBody(body),
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
