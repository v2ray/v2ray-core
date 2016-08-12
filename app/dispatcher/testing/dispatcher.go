package testing

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type TestPacketDispatcher struct {
	Destination chan v2net.Destination
	Handler     func(destination v2net.Destination, traffic ray.OutboundRay)
}

func NewTestPacketDispatcher(handler func(destination v2net.Destination, traffic ray.OutboundRay)) *TestPacketDispatcher {
	if handler == nil {
		handler = func(destination v2net.Destination, traffic ray.OutboundRay) {
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

func (this *TestPacketDispatcher) DispatchToOutbound(meta *proxy.InboundHandlerMeta, destination v2net.Destination) ray.InboundRay {
	traffic := ray.NewRay()
	this.Destination <- destination
	go this.Handler(destination, traffic)

	return traffic
}
