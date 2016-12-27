package dispatcher

import (
	"v2ray.com/core/app"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

const (
	APP_ID = app.ID(1)
)

// PacketDispatcher dispatch a packet and possibly further network payload to its destination.
type PacketDispatcher interface {
	DispatchToOutbound(session *proxy.SessionInfo) ray.InboundRay
}

type Inspector interface {
}
