package freedom

import (
	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/net"
)

type FreedomFactory struct {
}

func (factory FreedomFactory) Create(vp *core.Point, config []byte, dest v2net.Address) (core.OutboundConnectionHandler, error) {
	return NewFreedomConnection(dest), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
