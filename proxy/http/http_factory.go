package http

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

type HttpProxyServerFactory struct {
}

func (this HttpProxyServerFactory) Create(space app.Space, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
	return NewHttpProxyServer(space, rawConfig.(Config)), nil
}

func init() {
	connhandler.RegisterInboundConnectionHandlerFactory("http", HttpProxyServerFactory{})
}
