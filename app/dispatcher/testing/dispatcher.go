package testing

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type TestPacketDispatcher struct {
	LastPacket v2net.Packet
	Handler    func(packet v2net.Packet, traffic ray.OutboundRay)
}

func (this *TestPacketDispatcher) DispatchToOutbound(packet v2net.Packet) ray.InboundRay {
	traffic := ray.NewRay()
	this.LastPacket = packet
	if this.Handler == nil {
		go func() {
			for payload := range traffic.OutboundInput() {
				traffic.OutboundOutput() <- payload.Prepend([]byte("Processed: "))
			}
			close(traffic.OutboundOutput())
		}()
	} else {
		go this.Handler(packet, traffic)
	}

	return traffic
}
