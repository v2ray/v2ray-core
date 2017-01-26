// Package proxy contains all proxies used by V2Ray.
package proxy

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// An InboundHandler handles inbound network connections to V2Ray.
type InboundHandler interface {
	Network() net.NetworkList

	Process(context.Context, net.Network, internet.Connection) error
}

// An OutboundHandler handles outbound network connection for V2Ray.
type OutboundHandler interface {
	Process(context.Context, ray.OutboundRay) error
}

// Dialer is used by OutboundHandler for creating outbound connections.
type Dialer interface {
	Dial(ctx context.Context, destination net.Destination) (internet.Connection, error)
}
