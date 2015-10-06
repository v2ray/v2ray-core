package freedom

import (
	"github.com/v2ray/v2ray-core"
)

type FreedomFactory struct {
}

func (factory FreedomFactory) Create(vp *core.Point, config interface{}) (core.OutboundConnectionHandler, error) {
	return NewFreedomConnection(), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
