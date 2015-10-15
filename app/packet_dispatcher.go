package app

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// PacketDispatcher dispatch a packet and possibly further network payload to
// its destination.
type PacketDispatcher interface {
	DispatchToOutbound(packet v2net.Packet) ray.InboundRay
}
