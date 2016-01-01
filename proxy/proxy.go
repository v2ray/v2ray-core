// Package proxy contains all proxies used by V2Ray.

package proxy

import (
	"errors"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

var (
	inboundFactories  = make(map[string]internal.InboundConnectionHandlerCreator)
	outboundFactories = make(map[string]internal.OutboundConnectionHandlerCreator)

	ErrorProxyNotFound = errors.New("Proxy not found.")
	ErrorNameExists    = errors.New("Proxy with the same name already exists.")
)

func RegisterInboundConnectionHandlerFactory(name string, creator internal.InboundConnectionHandlerCreator) error {
	if _, found := inboundFactories[name]; found {
		return ErrorNameExists
	}
	inboundFactories[name] = creator
	return nil
}

func RegisterOutboundConnectionHandlerFactory(name string, creator internal.OutboundConnectionHandlerCreator) error {
	if _, found := outboundFactories[name]; found {
		return ErrorNameExists
	}
	outboundFactories[name] = creator
	return nil
}

func CreateInboundConnectionHandler(name string, space app.Space, config interface{}) (connhandler.InboundConnectionHandler, error) {
	if creator, found := inboundFactories[name]; !found {
		return nil, ErrorProxyNotFound
	} else {
		return creator(space, config)
	}
}

func CreateOutboundConnectionHandler(name string, space app.Space, config interface{}) (connhandler.OutboundConnectionHandler, error) {
	if creator, found := outboundFactories[name]; !found {
		return nil, ErrorNameExists
	} else {
		return creator(space, config)
	}
}
