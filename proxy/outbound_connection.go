package proxy

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type OutboundConnectionHandlerFactory interface {
	Create(config interface{}) (OutboundConnectionHandler, error)
}

type OutboundConnectionHandler interface {
	Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error
}
