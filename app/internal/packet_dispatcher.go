package internal

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type PacketDispatcherWithContext interface {
	DispatchToOutbound(context app.Context, packet v2net.Packet) ray.InboundRay
}

type contextedPacketDispatcher struct {
	context          app.Context
	packetDispatcher PacketDispatcherWithContext
}

func (this *contextedPacketDispatcher) DispatchToOutbound(packet v2net.Packet) ray.InboundRay {
	return this.packetDispatcher.DispatchToOutbound(this.context, packet)
}
