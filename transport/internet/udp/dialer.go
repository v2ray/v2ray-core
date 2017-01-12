package udp

import (
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
)

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_UDP,
		func(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
			conn, err := internet.DialSystem(src, dest)
			if err != nil {
				return nil, err
			}
			// TODO: handle dialer options
			return internal.NewConnection(internal.NewConnectionID(src, dest), conn, internal.NoOpConnectionRecyler{}, internal.ReuseConnection(false)), nil
		}))
}
