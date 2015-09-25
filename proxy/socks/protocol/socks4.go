package protocol

import (
	_ "fmt"

	"github.com/v2ray/v2ray-core/common/errors"
	_ "github.com/v2ray/v2ray-core/common/log"
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

func (r *Socks4AuthenticationResponse) ToBytes(buffer []byte) []byte {
	if buffer == nil {
		buffer = make([]byte, 8)
	}
	buffer[1] = r.result
	buffer[2] = byte(r.port >> 8)
	buffer[3] = byte(r.port)
	copy(buffer[4:], r.ip)
	return buffer
}
