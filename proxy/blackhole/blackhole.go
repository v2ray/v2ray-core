package blackhole

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/registry"
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

func (this *BlackHole) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	payload.Release()

	this.response.WriteTo(ray.OutboundOutput())
	ray.OutboundOutput().Close()

	ray.OutboundInput().Release()

	return nil
}

type Factory struct{}

func (this *Factory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_RawTCP},
	}
}

func (this *Factory) Create(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	return NewBlackHole(space, config.(*Config), meta)
}

func init() {
	registry.MustRegisterOutboundHandlerCreator(loader.GetType(new(Config)), new(Factory))
}
