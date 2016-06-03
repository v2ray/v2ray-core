package blackhole

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
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

func (this *BlackHole) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	payload.Release()

	ray.OutboundOutput().Close()
	ray.OutboundOutput().Release()

	ray.OutboundInput().Close()
	ray.OutboundInput().Release()

	return nil
}

func init() {
	internal.MustRegisterOutboundHandlerCreator("blackhole",
		func(space app.Space, config interface{}, sendThrough v2net.Address) (proxy.OutboundHandler, error) {
			return NewBlackHole(), nil
		})
}
