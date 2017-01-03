package udp

import (
	"net"

	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

type Connection struct {
	net.UDPConn
}

func (v *Connection) Reusable() bool {
	return false
}

func (v *Connection) SetReusable(b bool) {}

func init() {
	common.Must(internet.RegisterNetworkDialer(v2net.Network_UDP,
		func(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
			conn, err := internet.DialSystem(src, dest)
			if err != nil {
				return nil, err
			}
			// TODO: handle dialer options
			return &Connection{
				UDPConn: *(conn.(*net.UDPConn)),
			}, nil
		}))
}
