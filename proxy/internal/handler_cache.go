package internal

import (
	"errors"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

var (
	inboundFactories  = make(map[string]InboundHandlerCreator)
	outboundFactories = make(map[string]OutboundHandlerCreator)

	ErrorProxyNotFound    = errors.New("Proxy not found.")
	ErrorNameExists       = errors.New("Proxy with the same name already exists.")
	ErrorBadConfiguration = errors.New("Bad proxy configuration.")
)

func RegisterInboundHandlerCreator(name string, creator InboundHandlerCreator) error {
	if _, found := inboundFactories[name]; found {
		return ErrorNameExists
	}
	inboundFactories[name] = creator
	return nil
}

func MustRegisterInboundHandlerCreator(name string, creator InboundHandlerCreator) {
	if err := RegisterInboundHandlerCreator(name, creator); err != nil {
		panic(err)
	}
}

func RegisterOutboundHandlerCreator(name string, creator OutboundHandlerCreator) error {
	if _, found := outboundFactories[name]; found {
		return ErrorNameExists
	}
	outboundFactories[name] = creator
	return nil
}

func MustRegisterOutboundHandlerCreator(name string, creator OutboundHandlerCreator) {
	if err := RegisterOutboundHandlerCreator(name, creator); err != nil {
		panic(err)
	}
}

func CreateInboundHandler(name string, space app.Space, rawConfig []byte) (proxy.InboundHandler, error) {
	creator, found := inboundFactories[name]
	if !found {
		return nil, ErrorProxyNotFound
	}
	if len(rawConfig) > 0 {
		proxyConfig, err := config.CreateInboundConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator(space, proxyConfig)
	}
	return creator(space, nil)
}

func CreateOutboundHandler(name string, space app.Space, rawConfig []byte) (proxy.OutboundHandler, error) {
	creator, found := outboundFactories[name]
	if !found {
		return nil, ErrorProxyNotFound
	}

	if len(rawConfig) > 0 {
		proxyConfig, err := config.CreateOutboundConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator(space, proxyConfig)
	}

	return creator(space, nil)
}
