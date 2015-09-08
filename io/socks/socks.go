// Package socks contains protocol definition and io lib for SOCKS5 protocol
package socks

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	socksVersion = uint8(5)
)

// Authentication request header of Socks5 protocol
type Socks5AuthenticationRequest struct {
	version     byte
	nMethods    byte
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

	auth.nMethods = buffer[1]
	if auth.nMethods <= 0 {
		err = fmt.Errorf("Zero length of authentication methods")
		return
	}

	buffer = make([]byte, auth.nMethods)
	nBytes, err = reader.Read(buffer)
	if err != nil {
		return
	}
	if nBytes != int(auth.nMethods) {
		err = fmt.Errorf("Unmatching number of auth methods, expecting %d, but got %d", auth.nMethods, nBytes)
		return
	}
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
	return err
}

const (
	AddrTypeIPv4   = byte(0x01)
	AddrTypeIPv6   = byte(0x04)
	AddrTypeDomain = byte(0x03)
)

type Socks5Request struct {
	version  byte
	command  byte
	addrType byte
	ipv4     [4]byte
	domain   string
	ipv6     [16]byte
	port     uint16
}

func ReadRequest(reader io.Reader) (request *Socks5Request, err error) {
	request = new(Socks5Request)
	buffer := make([]byte, 4)
	nBytes, err := reader.Read(buffer)
	if err != nil {
		return
	}
	if nBytes < len(buffer) {
		err = fmt.Errorf("Unable to read request.")
		return
	}

	request.version = buffer[0]
	request.command = buffer[1]
	// buffer[2] is a reserved field
	request.addrType = buffer[3]
	switch request.addrType {
	case 0x01:
		nBytes, err = reader.Read(request.ipv4[:])
		if err != nil {
			return
		}
		if nBytes != 4 {
			err = fmt.Errorf("Unable to read IPv4 address.")
			return
		}
	case 0x03:
		buffer = make([]byte, 257)
		nBytes, err = reader.Read(buffer)
		if err != nil {
			return
		}
		domainLength := buffer[0]
		if nBytes != int(domainLength)+1 {
			err = fmt.Errorf("Unable to read domain")
			return
		}
		request.domain = string(buffer[1 : domainLength+1])
	case 0x04:
		nBytes, err = reader.Read(request.ipv6[:])
		if err != nil {
			return
		}
		if nBytes != 16 {
			err = fmt.Errorf("Unable to read IPv4 address.")
			return
		}
	default:
		err = fmt.Errorf("Unexpected address type %d", request.addrType)
		return
	}

	buffer = make([]byte, 2)
	nBytes, err = reader.Read(buffer)
	if err != nil {
		return
	}
	if nBytes != 2 {
		err = fmt.Errorf("Unable to read port.")
		return
	}

	request.port = binary.BigEndian.Uint16(buffer)
	return
}

const (
	ErrorSuccess                 = byte(0x00)
	ErrorGeneralFailure          = byte(0x01)
	ErrorConnectionNotAllowed    = byte(0x02)
	ErrorNetworkUnreachable      = byte(0x03)
	ErrorHostUnUnreachable       = byte(0x04)
	ErrorConnectionRefused       = byte(0x05)
	ErrorTTLExpired              = byte(0x06)
	ErrorCommandNotSupported     = byte(0x07)
	ErrorAddressTypeNotSupported = byte(0x08)
)

type Socks5Response struct {
	Version  byte
	Error    byte
	AddrType byte
	IPv4     [4]byte
	Domain   string
	IPv6     [16]byte
	Port     uint16
}

func (r Socks5Response) toBytes() []byte {
	buffer := make([]byte, 0, 300)
	buffer = append(buffer, r.Version)
	buffer = append(buffer, r.Error)
	buffer = append(buffer, 0x00) // reserved
	buffer = append(buffer, r.AddrType)
	switch r.AddrType {
	case 0x01:
		buffer = append(buffer, r.IPv4[:]...)
	case 0x03:
		buffer = append(buffer, byte(len(r.Domain)))
		buffer = append(buffer, []byte(r.Domain)...)
	case 0x04:
		buffer = append(buffer, r.IPv6[:]...)
	}
	portBuffer := make([]byte, 2)
	binary.BigEndian.PutUint16(portBuffer, r.Port)
	buffer = append(buffer, portBuffer...)
	return buffer
}

func WriteResponse(writer io.Writer, response Socks5Response) error {
	_, err := writer.Write(response.toBytes())
	return err
}
