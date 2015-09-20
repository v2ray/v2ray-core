package freedom

import (
	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type FreedomFactory struct {
}

func (factory FreedomFactory) Create(vp *core.Point, config []byte, dest v2net.Destination) (core.OutboundConnectionHandler, error) {
	return NewFreedomConnection(dest), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
