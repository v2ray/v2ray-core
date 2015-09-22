package freedom

import (
	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type FreedomFactory struct {
}

func (factory FreedomFactory) Create(vp *core.Point, config []byte, firstPacket v2net.Packet) (core.OutboundConnectionHandler, error) {
	return NewFreedomConnection(firstPacket), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
