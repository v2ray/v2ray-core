package wildcard_router

import (
	"github.com/v2ray/v2ray-core/app/point/config"
	"github.com/v2ray/v2ray-core/app/router"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type WildcardRouter struct {
}

func (router *WildcardRouter) TakeDetour(packet v2net.Packet) (config.DetourTag, error) {
	return "", nil
}

type WildcardRouterFactory struct {
}

func (factory *WildcardRouterFactory) Create(rawConfig interface{}) (router.Router, error) {
	return &WildcardRouter{}, nil
}

func init() {
	router.RegisterRouter("wildcard", &WildcardRouterFactory{})
}
