package testing

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
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

func (this *TestPacketDispatcher) DispatchToOutbound(session *proxy.SessionInfo) ray.InboundRay {
	traffic := ray.NewRay()
	this.Destination <- session.Destination
	go this.Handler(session.Destination, traffic)

	return traffic
}
