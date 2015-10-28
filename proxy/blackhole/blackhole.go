package blackhole

import (
	"io/ioutil"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// BlackHole is an outbound connection that sliently swallow the entire payload.
type BlackHole struct {
}

func NewBlackHole() *BlackHole {
	return &BlackHole{}
}

func (bh *BlackHole) Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error {
	if chunk := firstPacket.Chunk(); chunk != nil {
		chunk.Release()
	}

	close(ray.OutboundOutput())
	if firstPacket.MoreChunks() {
		v2net.ChanToWriter(ioutil.Discard, ray.OutboundInput())
	}
	return nil
}

type BlackHoleFactory struct {
}

func (factory BlackHoleFactory) Create(config interface{}) (proxy.OutboundConnectionHandler, error) {
	return NewBlackHole(), nil
}

func init() {
	proxy.RegisterOutboundConnectionHandlerFactory("blackhole", BlackHoleFactory{})
}
