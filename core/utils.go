package core

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
)

func head2binary(head SendPacketHead) (buffer *bytes.Buffer, err error) {
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
		err := binary.Read(r, ProtocolEndian, d)
		if err != nil {
			panic(err)
		}
	}
}
func MustBinaryWrite(r io.Writer, data ...interface{}) {
	for _, d := range data {
		err := binary.Write(r, ProtocolEndian, d)
		if err != nil {
			panic(err)
		}
	}
}

func ParseSIDString(sid string) (userID uint32, session [16]byte, err error) {
	if len(sid) != 40 {
		err = errors.New("illegal sid length")
		return
	}
	userIDtmp, err := strconv.ParseUint(sid[:8], 16, 32)
	userID = uint32(userIDtmp)

	sessiontmp, err := hex.DecodeString(sid[8:40])
	if err != nil {
		return
	}
	copy(session[:], sessiontmp[:32])
	return
}
