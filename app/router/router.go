package router

import (
	"github.com/v2ray/v2ray-core/app"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Router interface {
	TakeDetour(v2net.Destination) (string, error)
}

type RouterFactory interface {
	Create(rawConfig interface{}, space app.Space) (Router, error)
}

var (
	routerCache = make(map[string]RouterFactory)
)

func RegisterRouter(name string, factory RouterFactory) error {
	// TODO: check name
	routerCache[name] = factory
	return nil
}

func CreateRouter(name string, rawConfig interface{}, space app.Space) (Router, error) {
	if factory, found := routerCache[name]; found {
		return factory.Create(rawConfig, space)
	}
	return nil, ErrorRouterNotFound
}
