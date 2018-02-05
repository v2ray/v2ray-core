// Package proxy contains all proxies used by V2Ray.
//
// To implement an inbound or outbound proxy, one needs to do the following:
// 1. Implement the interface(s) below.
// 2. Register a config creator through common.RegisterConfig.
package proxy

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg proxy -path Proxy

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// An Inbound processes inbound connections.
type Inbound interface {
	// Network returns a list of network that this inbound supports. Connections with not-supported networks will not be passed into Process().
	Network() net.NetworkList

	// Process processes a connection of given network. If necessary, the Inbound can dispatch the connection to an Outbound.
	Process(context.Context, net.Network, internet.Connection, core.Dispatcher) error
}

// An Outbound process outbound connections.
type Outbound interface {
	// Process processes the given connection. The given dialer may be used to dial a system outbound connection.
	Process(context.Context, ray.OutboundRay, Dialer) error
}

// Dialer is used by OutboundHandler for creating outbound connections.
type Dialer interface {
	// Dial dials a system connection to the given destination.
	Dial(ctx context.Context, destination net.Destination) (internet.Connection, error)
}

// UserManager is the interface for Inbounds and Outbounds that can manage their users.
type UserManager interface {
	// AddUser adds a new user.
	AddUser(context.Context, *protocol.User) error

	// RemoveUser removes an user by email.
	RemoveUser(context.Context, string) error
}

type GetInbound interface {
	GetInbound() Inbound
}

type GetOutbound interface {
	GetOutbound() Outbound
}