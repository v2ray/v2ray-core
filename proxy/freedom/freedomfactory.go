package freedom

import (
	"github.com/v2ray/v2ray-core/proxy"
)

type FreedomFactory struct {
}

func (factory FreedomFactory) Create(config interface{}) (proxy.OutboundConnectionHandler, error) {
	return NewFreedomConnection(), nil
}

func init() {
	proxy.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
