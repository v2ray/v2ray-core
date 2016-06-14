package internal

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/internet"
)

type InboundHandlerFactory interface {
	StreamCapability() internet.StreamConnectionType
	Create(space app.Space, config interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error)
}

type OutboundHandlerFactory interface {
	StreamCapability() internet.StreamConnectionType
	Create(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error)
}
