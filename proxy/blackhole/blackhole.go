package blackhole

import (
	"io/ioutil"

	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// BlackHole is an outbound connection that sliently swallow the entire payload.
type BlackHole struct {
}

func NewBlackHole() *BlackHole {
	return &BlackHole{}
}

func (this *BlackHole) Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error {
	if chunk := firstPacket.Chunk(); chunk != nil {
		chunk.Release()
	}

	close(ray.OutboundOutput())
	if firstPacket.MoreChunks() {
		v2net.ChanToWriter(ioutil.Discard, ray.OutboundInput())
	}
	return nil
}

func init() {
	if err := internal.RegisterOutboundConnectionHandlerFactory("blackhole", func(space app.Space, config interface{}) (proxy.OutboundConnectionHandler, error) {
		return NewBlackHole(), nil
	}); err != nil {
		panic(err)
	}
}
