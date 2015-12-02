package protocol

import (
	"errors"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	Socks4Downgrade = errors.New("Downgraded to Socks 4.")
)

type Socks4AuthenticationRequest struct {
	Version byte
	Command byte
	Port    v2net.Port
	IP      [4]byte
}

type Socks4AuthenticationResponse struct {
	result byte
	port   uint16
	ip     []byte
}

func NewSocks4AuthenticationResponse(result byte, port v2net.Port, ip []byte) *Socks4AuthenticationResponse {
	return &Socks4AuthenticationResponse{
		result: result,
		port:   port.Value(),
		ip:     ip,
	}
}

func (r *Socks4AuthenticationResponse) Write(buffer *alloc.Buffer) {
	buffer.AppendBytes(
		byte(0x00), r.result, byte(r.port>>8), byte(r.port),
		r.ip[0], r.ip[1], r.ip[2], r.ip[3])
}
