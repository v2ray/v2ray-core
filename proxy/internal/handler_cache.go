package internal

import (
	"errors"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

var (
	inboundFactories  = make(map[string]InboundConnectionHandlerCreator)
	outboundFactories = make(map[string]OutboundConnectionHandlerCreator)

	ErrorProxyNotFound    = errors.New("Proxy not found.")
	ErrorNameExists       = errors.New("Proxy with the same name already exists.")
	ErrorBadConfiguration = errors.New("Bad proxy configuration.")
)

func RegisterInboundConnectionHandlerFactory(name string, creator InboundConnectionHandlerCreator) error {
	if _, found := inboundFactories[name]; found {
		return ErrorNameExists
	}
	inboundFactories[name] = creator
	return nil
}

func MustRegisterInboundConnectionHandlerCreator(name string, creator InboundConnectionHandlerCreator) {
	if err := RegisterInboundConnectionHandlerFactory(name, creator); err != nil {
		panic(err)
	}
}

func RegisterOutboundConnectionHandlerFactory(name string, creator OutboundConnectionHandlerCreator) error {
	if _, found := outboundFactories[name]; found {
		return ErrorNameExists
	}
	outboundFactories[name] = creator
	return nil
}

func MustRegisterOutboundConnectionHandlerCreator(name string, creator OutboundConnectionHandlerCreator) {
	if err := RegisterOutboundConnectionHandlerFactory(name, creator); err != nil {
		panic(err)
	}
}

func CreateInboundConnectionHandler(name string, space app.Space, rawConfig []byte) (proxy.InboundConnectionHandler, error) {
	creator, found := inboundFactories[name]
	if !found {
		return nil, ErrorProxyNotFound
	}
	if len(rawConfig) > 0 {
		proxyConfig, err := config.CreateInboundConnectionConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator(space, proxyConfig)
	}
	return creator(space, nil)
}

func CreateOutboundConnectionHandler(name string, space app.Space, rawConfig []byte) (proxy.OutboundConnectionHandler, error) {
	creator, found := outboundFactories[name]
	if !found {
		return nil, ErrorNameExists
	}

	if len(rawConfig) > 0 {
		proxyConfig, err := config.CreateOutboundConnectionConfig(name, rawConfig)
		if err != nil {
			return nil, err
		}
		return creator(space, proxyConfig)
	}

	return creator(space, nil)
}
