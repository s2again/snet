package core

import (
	"bytes"
	"encoding/binary"
	"io"
)

func head2binary(head PacketHead) (buffer *bytes.Buffer, err error) {
	defer func() {
		if x := recover(); x != nil {
			err = x.(error)
			return
		}
	}()
	buffer = new(bytes.Buffer)
	MustBinaryWrite(buffer, head.length)
	MustBinaryWrite(buffer, head.version)
	MustBinaryWrite(buffer, head.command)
	MustBinaryWrite(buffer, head.userID)
	MustBinaryWrite(buffer, head.sequence)
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
		MustBinaryWrite(buffer, v)
	}
	return
}

func MustBinaryRead(r io.Reader, data ...interface{}) {
	for _, d := range data {
		err := binary.Read(r, binary.BigEndian, d)
		if err != nil {
			panic(err)
		}
	}
}
func MustBinaryWrite(r io.Writer, data ...interface{}) {
	for _, d := range data {
		err := binary.Write(r, binary.BigEndian, d)
		if err != nil {
			panic(err)
		}
	}
}
