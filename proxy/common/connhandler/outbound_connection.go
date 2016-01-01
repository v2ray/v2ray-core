package connhandler

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// An OutboundConnectionHandler handles outbound network connection for V2Ray.
type OutboundConnectionHandler interface {
	// Dispatch sends one or more Packets to its destination.
	Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error
}
