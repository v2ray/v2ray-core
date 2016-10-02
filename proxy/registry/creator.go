package registry

import (
	"v2ray.com/core/app"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type InboundHandlerFactory interface {
	StreamCapability() v2net.NetworkList
	Create(space app.Space, config interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error)
}

type OutboundHandlerFactory interface {
	StreamCapability() v2net.NetworkList
	Create(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error)
}
