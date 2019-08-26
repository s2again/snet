package connection

import (
	"bytes"
	"encoding/binary"
	"io"
)

func head2binary(head packetHead) (buffer *bytes.Buffer, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	buffer = new(bytes.Buffer)
	mustBinaryWrite(buffer, head.length)
	mustBinaryWrite(buffer, head.version)
	mustBinaryWrite(buffer, head.command)
	mustBinaryWrite(buffer, head.userID)
	mustBinaryWrite(buffer, head.sequence)
	return
}
func var2binary(values ...interface{}) (buffer *bytes.Buffer, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	buffer = new(bytes.Buffer)
	for _, v := range values {
		mustBinaryWrite(buffer, v)
	}
	return
}

func mustBinaryRead(r io.Reader, data ...interface{}) {
	for _, d := range data {
		err := binary.Read(r, binary.BigEndian, d)
		if err != nil {
			panic(err)
		}
	}
}
func mustBinaryWrite(r io.Writer, data ...interface{}) {
	for _, d := range data {
		err := binary.Write(r, binary.BigEndian, d)
		if err != nil {
			panic(err)
		}
	}
}
