package udp

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName,
		func(ctx context.Context, dest net.Destination) (internet.Connection, error) {
			conn, err := internet.DialSystem(ctx, dest)
			if err != nil {
				return nil, err
			}
			// TODO: handle dialer options
			return internet.Connection(conn), nil
		}))
}
