package router

import (
	"errors"

	"github.com/v2ray/v2ray-core/app/point/config"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	RouterNotFound = errors.New("Router not found.")
)

type Router interface {
	TakeDetour(v2net.Destination) (config.DetourTag, error)
}

type RouterFactory interface {
	Create(rawConfig interface{}) (Router, error)
}

var (
	routerCache = make(map[string]RouterFactory)
)

func RegisterRouter(name string, factory RouterFactory) error {
	// TODO: check name
	routerCache[name] = factory
	return nil
}

func CreateRouter(name string, rawConfig interface{}) (Router, error) {
	if factory, found := routerCache[name]; found {
		return factory.Create(rawConfig)
	}
	return nil, RouterNotFound
}
