package dispatcher

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

const (
	APP_ID = app.ID(1)
)

// PacketDispatcher dispatch a packet and possibly further network payload to its destination.
type PacketDispatcher interface {
	DispatchToOutbound(destination v2net.Destination) ray.InboundRay
}

type packetDispatcherWithContext interface {
	DispatchToOutbound(context app.Context, destination v2net.Destination) ray.InboundRay
}

type contextedPacketDispatcher struct {
	context          app.Context
	packetDispatcher packetDispatcherWithContext
}

func (this *contextedPacketDispatcher) DispatchToOutbound(destination v2net.Destination) ray.InboundRay {
	return this.packetDispatcher.DispatchToOutbound(this.context, destination)
}

func init() {
	app.Register(APP_ID, func(context app.Context, obj interface{}) interface{} {
		packetDispatcher := obj.(packetDispatcherWithContext)
		return &contextedPacketDispatcher{
			context:          context,
			packetDispatcher: packetDispatcher,
		}
	})
}
