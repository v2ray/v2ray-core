// Package socks contains protocol definition and io lib for SOCKS5 protocol
package socks

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/v2ray/v2ray-core/log"
	v2net "github.com/v2ray/v2ray-core/net"
)

const (
	socksVersion  = byte(0x05)
	socks4Version = byte(0x04)

	AuthNotRequired      = byte(0x00)
	AuthGssApi           = byte(0x01)
	AuthUserPass         = byte(0x02)
	AuthNoMatchingMethod = byte(0xFF)

	Socks4RequestGranted  = byte(90)
	Socks4RequestRejected = byte(91)
)

var (
	ErrorSocksVersion4 = errors.New("Using SOCKS version 4.")
)

// Authentication request header of Socks5 protocol
type Socks5AuthenticationRequest struct {
	version     byte
	nMethods    byte
	authMethods [256]byte
}

type Socks4AuthenticationRequest struct {
	Version byte
	Command byte
	Port    uint16
	IP      [4]byte
}

func (request *Socks5AuthenticationRequest) HasAuthMethod(method byte) bool {
	for i := 0; i < int(request.nMethods); i++ {
		if request.authMethods[i] == method {
			return true
		}
	}
	return false
}

func ReadAuthentication(reader io.Reader) (auth Socks5AuthenticationRequest, auth4 Socks4AuthenticationRequest, err error) {
	buffer := make([]byte, 256)
	nBytes, err := reader.Read(buffer)
	if err != nil {
		log.Error("Failed to read socks authentication: %v", err)
		return
	}
	if nBytes < 2 {
		err = fmt.Errorf("Expected 2 bytes read, but actaully %d bytes read", nBytes)
		return
	}

	if buffer[0] == socks4Version {
		auth4.Version = buffer[0]
		auth4.Command = buffer[1]
		auth4.Port = binary.BigEndian.Uint16(buffer[2:4])
		copy(auth4.IP[:], buffer[4:8])
		err = ErrorSocksVersion4
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

	if nBytes-2 != int(auth.nMethods) {
		err = fmt.Errorf("Unmatching number of auth methods, expecting %d, but got %d", auth.nMethods, nBytes)
		return
	}
	copy(auth.authMethods[:], buffer[2:nBytes])
	return
}

type Socks5AuthenticationResponse struct {
	version    byte
	authMethod byte
}

type Socks4AuthenticationResponse struct {
	result byte
	port   uint16
	ip     []byte
}

func NewAuthenticationResponse(authMethod byte) *Socks5AuthenticationResponse {
	return &Socks5AuthenticationResponse{
		version:    socksVersion,
		authMethod: authMethod,
	}
}

func NewSocks4AuthenticationResponse(result byte, port uint16, ip []byte) *Socks4AuthenticationResponse {
	return &Socks4AuthenticationResponse{
		result: result,
		port:   port,
		ip:     ip,
	}
}

func WriteAuthentication(writer io.Writer, r *Socks5AuthenticationResponse) error {
	_, err := writer.Write([]byte{r.version, r.authMethod})
	return err
}

func WriteSocks4AuthenticationResponse(writer io.Writer, r *Socks4AuthenticationResponse) error {
	buffer := make([]byte, 8)
	// buffer[0] is always 0
	buffer[1] = r.result
	binary.BigEndian.PutUint16(buffer[2:4], r.port)
	copy(buffer[4:], r.ip)
	_, err := writer.Write(buffer)
	return err
}

type Socks5UserPassRequest struct {
	version  byte
	username string
	password string
}

func (request Socks5UserPassRequest) IsValid(username string, password string) bool {
	return request.username == username && request.password == password
}

func ReadUserPassRequest(reader io.Reader) (request Socks5UserPassRequest, err error) {
	buffer := make([]byte, 256)
	_, err = reader.Read(buffer[0:2])
	if err != nil {
		return
	}
	request.version = buffer[0]
	nUsername := buffer[1]
	nBytes, err := reader.Read(buffer[:nUsername])
	if err != nil {
		return
	}
	request.username = string(buffer[:nBytes])

	_, err = reader.Read(buffer[0:1])
	if err != nil {
		return
	}
	nPassword := buffer[0]
	nBytes, err = reader.Read(buffer[:nPassword])
	if err != nil {
		return
	}
	request.password = string(buffer[:nBytes])
	return
}

type Socks5UserPassResponse struct {
	version byte
	status  byte
}

func NewSocks5UserPassResponse(status byte) Socks5UserPassResponse {
	return Socks5UserPassResponse{
		version: socksVersion,
		status:  status,
	}
}

func WriteUserPassResponse(writer io.Writer, response Socks5UserPassResponse) error {
	_, err := writer.Write([]byte{response.version, response.status})
	return err
}

const (
	AddrTypeIPv4   = byte(0x01)
	AddrTypeIPv6   = byte(0x04)
	AddrTypeDomain = byte(0x03)

	CmdConnect      = byte(0x01)
	CmdBind         = byte(0x02)
	CmdUdpAssociate = byte(0x03)
)

type Socks5Request struct {
	Version  byte
	Command  byte
	AddrType byte
	IPv4     [4]byte
	Domain   string
	IPv6     [16]byte
	Port     uint16
}

func ReadRequest(reader io.Reader) (request *Socks5Request, err error) {

	buffer := make([]byte, 4)
	nBytes, err := reader.Read(buffer)
	if err != nil {
		return
	}
	if nBytes < len(buffer) {
		err = fmt.Errorf("Unable to read request.")
		return
	}
	request = &Socks5Request{
		Version: buffer[0],
		Command: buffer[1],
		// buffer[2] is a reserved field
		AddrType: buffer[3],
	}
	switch request.AddrType {
	case AddrTypeIPv4:
		nBytes, err = reader.Read(request.IPv4[:])
		if err != nil {
			return
		}
		if nBytes != 4 {
			err = fmt.Errorf("Unable to read IPv4 address.")
			return
		}
	case AddrTypeDomain:
		buffer = make([]byte, 256)
		nBytes, err = reader.Read(buffer[0:1])
		if err != nil {
			return
		}
		domainLength := buffer[0]
		nBytes, err = reader.Read(buffer[:domainLength])
		if err != nil {
			return
		}

		if nBytes != int(domainLength) {
			err = fmt.Errorf("Unable to read domain with %d bytes, expecting %d bytes", nBytes, domainLength)
			return
		}
		request.Domain = string(buffer[:domainLength])
	case AddrTypeIPv6:
		nBytes, err = reader.Read(request.IPv6[:])
		if err != nil {
			return
		}
		if nBytes != 16 {
			err = fmt.Errorf("Unable to read IPv4 address.")
			return
		}
	default:
		err = fmt.Errorf("Unexpected address type %d", request.AddrType)
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

	request.Port = binary.BigEndian.Uint16(buffer)
	return
}

func (request *Socks5Request) Destination() v2net.Address {
	switch request.AddrType {
	case AddrTypeIPv4:
		return v2net.IPAddress(request.IPv4[:], request.Port)
	case AddrTypeIPv6:
		return v2net.IPAddress(request.IPv6[:], request.Port)
	case AddrTypeDomain:
		return v2net.DomainAddress(request.Domain, request.Port)
	default:
		panic("Unknown address type")
	}
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

func NewSocks5Response() *Socks5Response {
	return &Socks5Response{
		Version: socksVersion,
	}
}

func (r *Socks5Response) SetIPv4(ipv4 []byte) {
	r.AddrType = AddrTypeIPv4
	copy(r.IPv4[:], ipv4)
}

func (r *Socks5Response) SetIPv6(ipv6 []byte) {
	r.AddrType = AddrTypeIPv6
	copy(r.IPv6[:], ipv6)
}

func (r *Socks5Response) SetDomain(domain string) {
	r.AddrType = AddrTypeDomain
	r.Domain = domain
}

func (r *Socks5Response) toBytes() []byte {
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

func WriteResponse(writer io.Writer, response *Socks5Response) error {
	_, err := writer.Write(response.toBytes())
	return err
}
