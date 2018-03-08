package udp

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_UDP,
		func(ctx context.Context, dest net.Destination) (internet.Connection, error) {
			src := internet.DialerSourceFromContext(ctx)
			conn, err := internet.DialSystem(ctx, src, dest)
			if err != nil {
				return nil, err
			}
			// TODO: handle dialer options
			return internet.Connection(conn), nil
		}))
}
