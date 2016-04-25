// Package proxy contains all proxies used by V2Ray.
package proxy // import "github.com/v2ray/v2ray-core/proxy"

import (
	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type HandlerState int

const (
	HandlerStateStopped = HandlerState(0)
	HandlerStateRunning = HandlerState(1)
)

// An InboundHandler handles inbound network connections to V2Ray.
type InboundHandler interface {
	// Listen starts a InboundHandler by listen on a specific port.
	Listen(port v2net.Port) error
	// Close stops the handler to accepting anymore inbound connections.
	Close()
	// Port returns the port that the handler is listening on.
	Port() v2net.Port
}

// An OutboundHandler handles outbound network connection for V2Ray.
type OutboundHandler interface {
	// Dispatch sends one or more Packets to its destination.
	Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error
}
