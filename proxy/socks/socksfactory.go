package socks

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func init() {
	internal.MustRegisterInboundConnectionHandlerCreator("socks",
		func(space app.Space, rawConfig interface{}) (proxy.InboundConnectionHandler, error) {
			return NewSocksServer(space, rawConfig.(Config)), nil
		})
}
