package router

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
)

const (
	APP_ID = app.ID(3)
)

type Router interface {
	common.Releasable
	TakeDetour(v2net.Destination) (string, error)
}

type RouterFactory interface {
	Create(rawConfig interface{}, space app.Space) (Router, error)
}

var (
	routerCache = make(map[string]RouterFactory)
)

func RegisterRouter(name string, factory RouterFactory) error {
	if _, found := routerCache[name]; found {
		return common.ErrDuplicatedName
	}
	routerCache[name] = factory
	return nil
}

func CreateRouter(name string, rawConfig interface{}, space app.Space) (Router, error) {
	if factory, found := routerCache[name]; found {
		return factory.Create(rawConfig, space)
	}
	log.Error("Router: not found: ", name)
	return nil, common.ErrObjectNotFound
}
