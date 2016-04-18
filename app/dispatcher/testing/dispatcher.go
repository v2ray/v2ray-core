package testing

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type TestPacketDispatcher struct {
	Destination chan v2net.Destination
	Handler     func(packet v2net.Packet, traffic ray.OutboundRay)
}

func NewTestPacketDispatcher(handler func(packet v2net.Packet, traffic ray.OutboundRay)) *TestPacketDispatcher {
	if handler == nil {
		handler = func(packet v2net.Packet, traffic ray.OutboundRay) {
			for {
				payload, err := traffic.OutboundInput().Read()
				if err != nil {
					break
				}
				traffic.OutboundOutput().Write(payload.Prepend([]byte("Processed: ")))
			}
			traffic.OutboundOutput().Close()
		}
	}
	return &TestPacketDispatcher{
		Destination: make(chan v2net.Destination),
		Handler:     handler,
	}
}

func (this *TestPacketDispatcher) DispatchToOutbound(packet v2net.Packet) ray.InboundRay {
	traffic := ray.NewRay()
	this.Destination <- packet.Destination()
	go this.Handler(packet, traffic)

	return traffic
}
