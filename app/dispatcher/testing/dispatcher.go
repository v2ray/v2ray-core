package testing

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type TestPacketDispatcher struct {
	LastPacket chan v2net.Packet
	Handler    func(packet v2net.Packet, traffic ray.OutboundRay)
}

func NewTestPacketDispatcher(handler func(packet v2net.Packet, traffic ray.OutboundRay)) *TestPacketDispatcher {
	if handler == nil {
		handler = func(packet v2net.Packet, traffic ray.OutboundRay) {
			for payload := range traffic.OutboundInput() {
				traffic.OutboundOutput() <- payload.Prepend([]byte("Processed: "))
			}
			close(traffic.OutboundOutput())
		}
	}
	return &TestPacketDispatcher{
		LastPacket: make(chan v2net.Packet, 16),
		Handler:    handler,
	}
}

func (this *TestPacketDispatcher) DispatchToOutbound(packet v2net.Packet) ray.InboundRay {
	traffic := ray.NewRay()
	this.LastPacket <- packet
	go this.Handler(packet, traffic)

	return traffic
}
