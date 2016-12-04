package registry

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/proxy"
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

func MustRegisterInboundHandlerCreator(name string, creator InboundHandlerFactory) {
	if err := RegisterInboundHandlerCreator(name, creator); err != nil {
		panic(err)
	}
}

func RegisterOutboundHandlerCreator(name string, creator OutboundHandlerFactory) error {
	if _, found := outboundFactories[name]; found {
		return common.ErrDuplicatedName
	}
	outboundFactories[name] = creator
	return nil
}

func MustRegisterOutboundHandlerCreator(name string, creator OutboundHandlerFactory) {
	if err := RegisterOutboundHandlerCreator(name, creator); err != nil {
		panic(err)
	}
}

func CreateInboundHandler(name string, space app.Space, config interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	creator, found := inboundFactories[name]
	if !found {
		return nil, errors.New("Proxy|Registry: Unknown inbound name: " + name)
	}
	if meta.StreamSettings == nil {
		meta.StreamSettings = &internet.StreamConfig{
			Network: creator.StreamCapability().Get(0),
		}
	} else {
		if !creator.StreamCapability().HasNetwork(meta.StreamSettings.Network) {
			return nil, errors.New("Proxy|Registry: Invalid network: " + meta.StreamSettings.Network.String())
		}
	}

	return creator.Create(space, config, meta)
}

func CreateOutboundHandler(name string, space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	creator, found := outboundFactories[name]
	if !found {
		return nil, errors.New("Proxy|Registry: Unknown outbound name: " + name)
	}
	if meta.StreamSettings == nil {
		meta.StreamSettings = &internet.StreamConfig{
			Network: creator.StreamCapability().Get(0),
		}
	} else {
		if !creator.StreamCapability().HasNetwork(meta.StreamSettings.Network) {
			return nil, errors.New("Proxy|Registry: Invalid network: " + meta.StreamSettings.Network.String())
		}
	}

	return creator.Create(space, config, meta)
}
