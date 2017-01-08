package proxy

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

var (
	inboundFactories  = make(map[string]InboundHandlerFactory)
	outboundFactories = make(map[string]OutboundHandlerFactory)
)

func RegisterInboundHandlerCreator(name string, creator InboundHandlerFactory) error {
	if _, found := inboundFactories[name]; found {
		return common.ErrDuplicatedName
	}
	inboundFactories[name] = creator
	return nil
}

func RegisterOutboundHandlerCreator(name string, creator OutboundHandlerFactory) error {
	if _, found := outboundFactories[name]; found {
		return common.ErrDuplicatedName
	}
	outboundFactories[name] = creator
	return nil
}

func CreateInboundHandler(name string, space app.Space, config interface{}, meta *InboundHandlerMeta) (InboundHandler, error) {
	creator, found := inboundFactories[name]
	if !found {
		return nil, errors.New("Proxy: Unknown inbound name: " + name)
	}
	if meta.StreamSettings == nil {
		meta.StreamSettings = &internet.StreamConfig{
			Network: creator.StreamCapability().Get(0),
		}
	} else if meta.StreamSettings.Network == v2net.Network_Unknown {
		meta.StreamSettings.Network = creator.StreamCapability().Get(0)
	} else {
		if !creator.StreamCapability().HasNetwork(meta.StreamSettings.Network) {
			return nil, errors.New("Proxy: Invalid network: " + meta.StreamSettings.Network.String())
		}
	}

	return creator.Create(space, config, meta)
}

func CreateOutboundHandler(name string, space app.Space, config interface{}, meta *OutboundHandlerMeta) (OutboundHandler, error) {
	creator, found := outboundFactories[name]
	if !found {
		return nil, errors.New("Proxy: Unknown outbound name: " + name)
	}
	if meta.StreamSettings == nil {
		meta.StreamSettings = &internet.StreamConfig{
			Network: creator.StreamCapability().Get(0),
		}
	} else if meta.StreamSettings.Network == v2net.Network_Unknown {
		meta.StreamSettings.Network = creator.StreamCapability().Get(0)
	} else {
		if !creator.StreamCapability().HasNetwork(meta.StreamSettings.Network) {
			return nil, errors.New("Proxy: Invalid network: " + meta.StreamSettings.Network.String())
		}
	}

	return creator.Create(space, config, meta)
}
