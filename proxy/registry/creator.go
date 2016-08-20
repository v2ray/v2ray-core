package registry

import (
	"v2ray.com/core/app"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
)

type InboundHandlerFactory interface {
	StreamCapability() internet.StreamConnectionType
	Create(space app.Space, config interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error)
}

type OutboundHandlerFactory interface {
	StreamCapability() internet.StreamConnectionType
	Create(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error)
}
