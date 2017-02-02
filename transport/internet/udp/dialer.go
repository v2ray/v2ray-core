package udp

import (
	"context"

	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
)

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_UDP,
		func(ctx context.Context, dest v2net.Destination) (internet.Connection, error) {
			src := internet.DialerSourceFromContext(ctx)
			conn, err := internet.DialSystem(src, dest)
			if err != nil {
				return nil, err
			}
			// TODO: handle dialer options
			return internal.NewConnection(internal.NewConnectionID(src, dest), conn, internal.NoOpConnectionRecyler{}, internal.ReuseConnection(false)), nil
		}))
}
