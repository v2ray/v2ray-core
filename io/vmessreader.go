package io

import (
	"net"
)

// VMessInput implements the input message of VMess protocol. It only contains
// the header of a input message. The data part will be handled by conection
// handler directly, in favor of data streaming.
type VMessInput struct {
	version  byte
	userHash [16]byte
	randHash [256]byte
	respKey  [32]byte
	iv       [16]byte
	command  byte
	port     uint16
	target   [256]byte
}

type VMessReader struct {
	conn *net.Conn
}

// NewVMessReader creates a new VMessReader
func NewVMessReader(conn *net.Conn) (VMessReader, error) {
	var reader VMessReader
	reader.conn = conn
	return reader, nil
}

// Read reads data from current connection and form a VMessInput if possible.
func (*VMessReader) Read() (VMessInput, error) {
	var input VMessInput
	return input, nil
}
