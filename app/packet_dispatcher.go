package app

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// PacketDispatcher dispatch a packet and possibly further network payload to its destination.
type PacketDispatcher interface {
	DispatchToOutbound(packet v2net.Packet) ray.InboundRay
}

type PacketDispatcherWithContext interface {
	DispatchToOutbound(context Context, packet v2net.Packet) ray.InboundRay
}

type contextedPacketDispatcher struct {
	context          Context
	packetDispatcher PacketDispatcherWithContext
}

func (this *contextedPacketDispatcher) DispatchToOutbound(packet v2net.Packet) ray.InboundRay {
	return this.packetDispatcher.DispatchToOutbound(this.context, packet)
}
