package dispatcher

import (
	"v2ray.com/core/app"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

// PacketDispatcher dispatch a packet and possibly further network payload to its destination.
type PacketDispatcher interface {
	DispatchToOutbound(session *proxy.SessionInfo) ray.InboundRay
}

func FromSpace(space app.Space) PacketDispatcher {
	if app := space.GetApplication((*PacketDispatcher)(nil)); app != nil {
		return app.(PacketDispatcher)
	}
	return nil
}
