// Package proxy contains all proxies used by V2Ray.

package proxy

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// A InboundConnectionHandler handles inbound network connections to V2Ray.
type InboundConnectionHandler interface {
	// Listen starts a InboundConnectionHandler by listen on a specific port.
	Listen(port v2net.Port) error
	// Close stops the handler to accepting anymore inbound connections.
	Close()
	// Port returns the port that the handler is listening on.
	Port() v2net.Port
}

// An OutboundConnectionHandler handles outbound network connection for V2Ray.
type OutboundConnectionHandler interface {
	// Dispatch sends one or more Packets to its destination.
	Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error
}
