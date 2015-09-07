package io

import (
	"net"
)

// VMessInput implements the input message of VMess protocol.
type VMessInput struct {
	version  byte
	userHash [16]byte
	randHash [256]byte
	respKey  [32]byte
	iv       [16]byte
	command  byte
	port     uint16
	target   [256]byte
	data     []byte
}

type VMessReader struct {
	conn *net.Conn
}

func NewVMessReader(conn *net.Conn) (VMessReader, error) {
	var reader VMessReader
	reader.conn = conn
	return reader, nil
}

func (*VMessReader) Read() (VMessInput, error) {
	var input VMessInput
	return input, nil
}
