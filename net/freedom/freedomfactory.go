package freedom

import (
	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/net"
)

type FreedomFactory struct {
}

func (factory FreedomFactory) Create(vp *core.VPoint, config []byte, dest v2net.VAddress) (core.OutboundConnectionHandler, error) {
	return NewVFreeConnection(dest), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
