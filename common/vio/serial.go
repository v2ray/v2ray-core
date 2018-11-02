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

func ReadUint16(reader io.Reader) (uint16, error) {
	var b [2]byte
	if _, err := io.ReadFull(reader, b[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(b[:]), nil
}

func WriteUint16(writer io.Writer, value uint16) (int, error) {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], value)
	return writer.Write(b[:])
}
