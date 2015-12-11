package connhandler

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// An OutboundConnectionHandlerFactory creates OutboundConnectionHandler on demand.
type OutboundConnectionHandlerFactory interface {
	// Create creates a new OutboundConnectionHandler with given config.
	Create(space app.Space, config interface{}) (OutboundConnectionHandler, error)
}

// An OutboundConnectionHandler handles outbound network connection for V2Ray.
type OutboundConnectionHandler interface {
	// Dispatch sends one or more Packets to its destination.
	Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error
}
