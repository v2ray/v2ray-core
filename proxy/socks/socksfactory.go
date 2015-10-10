package socks

import (
	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/proxy/socks/config/json"
)

type SocksServerFactory struct {
}

func (factory SocksServerFactory) Create(vp *core.Point, rawConfig interface{}) (core.InboundConnectionHandler, error) {
	config := rawConfig.(*json.SocksConfig)
	config.Initialize()
	return NewSocksServer(vp, config), nil
}

func init() {
	core.RegisterInboundConnectionHandlerFactory("socks", SocksServerFactory{})
}
