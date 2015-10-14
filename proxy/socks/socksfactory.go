package socks

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/socks/config/json"
)

type SocksServerFactory struct {
}

func (factory SocksServerFactory) Create(dispatcher app.PacketDispatcher, rawConfig interface{}) (proxy.InboundConnectionHandler, error) {
	config := rawConfig.(*json.SocksConfig)
	config.Initialize()
	return NewSocksServer(dispatcher, config), nil
}

func init() {
	proxy.RegisterInboundConnectionHandlerFactory("socks", SocksServerFactory{})
}
