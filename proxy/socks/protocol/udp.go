package protocol

import (
	"io"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Socks5UDPRequest struct {
	fragment byte
	address  v2net.Address
	data     []byte
}

func ReadUDPRequest(reader io.Reader) (request Socks5UDPRequest, err error) {
	//buf := make([]byte, 4 * 1024) // Regular UDP packet size is 1500 bytes.

	return
}
