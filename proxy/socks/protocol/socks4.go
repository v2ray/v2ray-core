package protocol

import (
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/errors"
)

type SocksVersion4Error struct {
	errors.ErrorCode
}

var socksVersion4ErrorInstance = SocksVersion4Error{ErrorCode: 1000}

func NewSocksVersion4Error() SocksVersion4Error {
	return socksVersion4ErrorInstance
}

func (err SocksVersion4Error) Error() string {
	return err.Prefix() + "Request is socks version 4."
}

type Socks4AuthenticationRequest struct {
	Version byte
	Command byte
	Port    uint16
	IP      [4]byte
}

type Socks4AuthenticationResponse struct {
	result byte
	port   uint16
	ip     []byte
}

func NewSocks4AuthenticationResponse(result byte, port uint16, ip []byte) *Socks4AuthenticationResponse {
	return &Socks4AuthenticationResponse{
		result: result,
		port:   port,
		ip:     ip,
	}
}

func (r *Socks4AuthenticationResponse) Write(buffer *alloc.Buffer) {
	buffer.AppendBytes(
		byte(0x00), r.result, byte(r.port>>8), byte(r.port),
		r.ip[0], r.ip[1], r.ip[2], r.ip[3])
}
