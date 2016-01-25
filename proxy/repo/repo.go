package repo

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func CreateInboundHandler(name string, space app.Space, rawConfig []byte) (proxy.InboundHandler, error) {
	return internal.CreateInboundHandler(name, space, rawConfig)
}

func CreateOutboundHandler(name string, space app.Space, rawConfig []byte) (proxy.OutboundHandler, error) {
	return internal.CreateOutboundHandler(name, space, rawConfig)
}
