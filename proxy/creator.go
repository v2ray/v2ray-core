package proxy

import (
	"v2ray.com/core/app"
	v2net "v2ray.com/core/common/net"
)

type InboundHandlerFactory interface {
	StreamCapability() v2net.NetworkList
	Create(space app.Space, config interface{}, meta *InboundHandlerMeta) (InboundHandler, error)
}

type OutboundHandlerFactory interface {
	StreamCapability() v2net.NetworkList
	Create(space app.Space, config interface{}, meta *OutboundHandlerMeta) (OutboundHandler, error)
}
