// Package proxy contains all proxies used by V2Ray.
package proxy

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type HandlerState int

const (
	HandlerStateStopped = HandlerState(0)
	HandlerStateRunning = HandlerState(1)
)

type SessionInfo struct {
	Source      net.Destination
	Destination net.Destination
	User        *protocol.User
	Inbound     *InboundHandlerMeta
}

type InboundHandlerMeta struct {
	Tag                    string
	Address                net.Address
	Port                   net.Port
	AllowPassiveConnection bool
	StreamSettings         *internet.StreamConfig
}

type OutboundHandlerMeta struct {
	Tag            string
	Address        net.Address
	StreamSettings *internet.StreamConfig
	ProxySettings  *internet.ProxyConfig
}

func (v *OutboundHandlerMeta) GetDialerOptions() internet.DialerOptions {
	return internet.DialerOptions{
		Stream: v.StreamSettings,
		Proxy:  v.ProxySettings,
	}
}

// An InboundHandler handles inbound network connections to V2Ray.
type InboundHandler interface {
	// Listen starts a InboundHandler.
	Start() error
	// Close stops the handler to accepting anymore inbound connections.
	Close()
	// Port returns the port that the handler is listening on.
	Port() net.Port
}

// An OutboundHandler handles outbound network connection for V2Ray.
type OutboundHandler interface {
	// Dispatch sends one or more Packets to its destination.
	Dispatch(destination net.Destination, ray ray.OutboundRay)
}

// Dialer is used by OutboundHandler for creating outbound connections.
type Dialer interface {
	Dial(ctx context.Context, destination net.Destination) (internet.Connection, error)
}
