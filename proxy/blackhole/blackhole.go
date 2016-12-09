package blackhole

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common/buf"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

// BlackHole is an outbound connection that sliently swallow the entire payload.
type BlackHole struct {
	meta     *proxy.OutboundHandlerMeta
	response ResponseConfig
}

func NewBlackHole(space app.Space, config *Config, meta *proxy.OutboundHandlerMeta) (*BlackHole, error) {
	response, err := config.GetInternalResponse()
	if err != nil {
		return nil, err
	}
	return &BlackHole{
		meta:     meta,
		response: response,
	}, nil
}

func (v *BlackHole) Dispatch(destination v2net.Destination, payload *buf.Buffer, ray ray.OutboundRay) {
	payload.Release()

	v.response.WriteTo(ray.OutboundOutput())
	ray.OutboundOutput().Close()

	ray.OutboundInput().Release()
}

type Factory struct{}

func (v *Factory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_RawTCP},
	}
}

func (v *Factory) Create(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	return NewBlackHole(space, config.(*Config), meta)
}
