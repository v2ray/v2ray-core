package udp

import (
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"
)

type Connection struct {
	net.UDPConn
}

func (this *Connection) Reusable() bool {
	return false
}

func (this *Connection) SetReusable(b bool) {}

func init() {
	internet.UDPDialer = func(src v2net.Address, dest v2net.Destination) (internet.Connection, error) {
		conn, err := internet.DialToDest(src, dest)
		if err != nil {
			return nil, err
		}
		return &Connection{
			UDPConn: *(conn.(*net.UDPConn)),
		}, nil
	}
}
