package internal

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
)

type InboundHandlerCreator func(space app.Space, config interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error)
type OutboundHandlerCreator func(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error)
