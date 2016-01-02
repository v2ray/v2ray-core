package json

import (
	"github.com/v2ray/v2ray-core/proxy/internal/config"
	"github.com/v2ray/v2ray-core/proxy/internal/config/json"
)

type FreedomConfiguration struct {
}

func init() {
	config.RegisterOutboundConnectionConfig("freedom", json.JsonConfigLoader(func() interface{} {
		return &FreedomConfiguration{}
	}))
}
