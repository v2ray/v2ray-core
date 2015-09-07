// Package socks contains protocol definition and io lib for SOCKS5 protocol
package socks

import (
	"fmt"
	"io"
)

const (
	socksVersion = uint8(5)
)

// Authentication request header of Socks5 protocol
type Socks5AuthenticationRequest struct {
	version     byte
	authMethods [256]byte
}

func ReadAuthentication(reader io.Reader) (auth Socks5AuthenticationRequest, err error) {
	buffer := make([]byte, 2)
	nBytes, err := reader.Read(buffer)
	if err != nil {
		return
	}
	if nBytes < 2 {
		err = fmt.Errorf("Expected 2 bytes read, but actaully %d bytes read", nBytes)
		return
	}

	auth.version = buffer[0]
	if auth.version != socksVersion {
		err = fmt.Errorf("Unknown SOCKS version %d", auth.version)
		return
	}

	nMethods := buffer[1]
	if nMethods <= 0 {
		err = fmt.Errorf("Zero length of authentication methods")
		return
	}

	buffer = make([]byte, nMethods)
	nBytes, err = reader.Read(buffer)
	copy(auth.authMethods[:nBytes], buffer)
	return
}

type Socks5AuthenticationResponse struct {
	version    byte
	authMethod byte
}

func (r *Socks5AuthenticationResponse) ToBytes() []byte {
	buffer := make([]byte, 2 /* size of Socks5AuthenticationResponse */)
	buffer[0] = r.version
	buffer[1] = r.authMethod
	return buffer
}

func WriteAuthentication(writer io.Writer, response Socks5AuthenticationResponse) error {
	_, err := writer.Write(response.ToBytes())
	if err != nil {
		return err
	}
	return nil
}
