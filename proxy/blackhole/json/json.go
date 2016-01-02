package json

import (
	"github.com/v2ray/v2ray-core/proxy/internal/config"
	"github.com/v2ray/v2ray-core/proxy/internal/config/json"
)

type BlackHoleConfig struct {
}

func init() {
	config.RegisterOutboundConnectionConfig("blackhole", json.JsonConfigLoader(func() interface{} {
		return new(BlackHoleConfig)
	}))
}
