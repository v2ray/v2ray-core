package internal

import (
	"errors"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/internet"
)

var (
	inboundFactories  = make(map[string]InboundHandlerFactory)
	outboundFactories = make(map[string]OutboundHandlerFactory)

	ErrorProxyNotFound    = errors.New("Proxy not found.")
	ErrorNameExists       = errors.New("Proxy with the same name already exists.")
	ErrorBadConfiguration = errors.New("Bad proxy configuration.")
)

func RegisterInboundHandlerCreator(name string, creator InboundHandlerFactory) error {
	if _, found := inboundFactories[name]; found {
		return ErrorNameExists
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
		return ErrorNameExists
	}
	outboundFactories[name] = creator
	return nil
}

func MustRegisterOutboundHandlerCreator(name string, creator OutboundHandlerFactory) {
	if err := RegisterOutboundHandlerCreator(name, creator); err != nil {
		panic(err)
	}
}

func CreateInboundHandler(name string, space app.Space, rawConfig []byte, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	creator, found := inboundFactories[name]
	if !found {
		return nil, ErrorProxyNotFound
	}
	if meta.StreamSettings == nil {
		meta.StreamSettings = &internet.StreamSettings{
			Type: creator.StreamCapability(),
		}
	} else {
		meta.StreamSettings.Type &= creator.StreamCapability()
	}

	if len(rawConfig) > 0 {
		proxyConfig, err := CreateInboundConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator.Create(space, proxyConfig, meta)
	}
	return creator.Create(space, nil, meta)
}

func CreateOutboundHandler(name string, space app.Space, rawConfig []byte, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	creator, found := outboundFactories[name]
	if !found {
		return nil, ErrorProxyNotFound
	}
	if meta.StreamSettings == nil {
		meta.StreamSettings = &internet.StreamSettings{
			Type: creator.StreamCapability(),
		}
	} else {
		meta.StreamSettings.Type &= creator.StreamCapability()
	}

	if len(rawConfig) > 0 {
		proxyConfig, err := CreateOutboundConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator.Create(space, proxyConfig, meta)
	}

	return creator.Create(space, nil, meta)
}
