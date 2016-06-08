package repo

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func CreateInboundHandler(name string, space app.Space, rawConfig []byte, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	return internal.CreateInboundHandler(name, space, rawConfig, meta)
}

func CreateOutboundHandler(name string, space app.Space, rawConfig []byte, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	return internal.CreateOutboundHandler(name, space, rawConfig, meta)
}
