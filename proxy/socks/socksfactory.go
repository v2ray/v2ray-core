package socks

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/socks/config"
)

type SocksServerFactory struct {
}

func (this SocksServerFactory) Create(space *app.Space, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
	return NewSocksServer(space, rawConfig.(config.SocksConfig)), nil
}

func init() {
	connhandler.RegisterInboundConnectionHandlerFactory("socks", SocksServerFactory{})
}
