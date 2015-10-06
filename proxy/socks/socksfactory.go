package socks

import (
	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/proxy/socks/config/json"
)

type SocksServerFactory struct {
}

func (factory SocksServerFactory) Create(vp *core.Point, config interface{}) (core.InboundConnectionHandler, error) {
	return NewSocksServer(vp, config.(*json.SocksConfig)), nil
}

func init() {
	core.RegisterInboundConnectionHandlerFactory("socks", SocksServerFactory{})
}
