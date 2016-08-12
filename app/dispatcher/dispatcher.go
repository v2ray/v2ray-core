package dispatcher

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/ray"
)

const (
	APP_ID = app.ID(1)
)

// PacketDispatcher dispatch a packet and possibly further network payload to its destination.
type PacketDispatcher interface {
	DispatchToOutbound(meta *proxy.InboundHandlerMeta, destination v2net.Destination) ray.InboundRay
}
