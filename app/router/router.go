package router

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/config"
)

type Router interface {
	TakeDetour(v2net.Packet) (config.ConnectionTag, error)
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
