// Package proxy contains all proxies used by V2Ray.
package proxy

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// An Inbound processes inbound connections.
type Inbound interface {
	Network() net.NetworkList

	Process(context.Context, net.Network, internet.Connection) error
}

// An Outbound process outbound connections.
type Outbound interface {
	Process(context.Context, ray.OutboundRay) error
}

// Dialer is used by OutboundHandler for creating outbound connections.
type Dialer interface {
	Dial(ctx context.Context, destination net.Destination) (internet.Connection, error)
}
