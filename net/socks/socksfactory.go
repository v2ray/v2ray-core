package socks

import (
	"github.com/v2ray/v2ray-core"
)

type SocksServerFactory struct {
}

func (factory SocksServerFactory) Create(vp *core.Point, config []byte) (core.InboundConnectionHandler, error) {
	return NewSocksServer(vp, config), nil
}

func init() {
	core.RegisterInboundConnectionHandlerFactory("socks", SocksServerFactory{})
}
