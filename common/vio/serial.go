package vio

import (
	"encoding/binary"
	"io"
)

func WriteUint32(writer io.Writer, value uint32) (int, error) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], value)
	return writer.Write(b[:])
}

func WriteUint16(writer io.Writer, value uint16) (int, error) {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], value)
	return writer.Write(b[:])
}
