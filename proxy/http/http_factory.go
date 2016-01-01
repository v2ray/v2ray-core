package http

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

func init() {
	if err := proxy.RegisterInboundConnectionHandlerFactory("http", func(space app.Space, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
		return NewHttpProxyServer(space, rawConfig.(Config)), nil
	}); err != nil {
		panic(err)
	}
}
