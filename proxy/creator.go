package proxy

import "v2ray.com/core/app"

type InboundHandlerFactory interface {
	Create(space app.Space, config interface{}, meta *InboundHandlerMeta) (InboundHandler, error)
}

type OutboundHandlerFactory interface {
	Create(space app.Space, config interface{}, meta *OutboundHandlerMeta) (OutboundHandler, error)
}
