package app

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type PacketDispatcher interface {
	DispatchToOutbound(packet v2net.Packet) ray.InboundRay
}
