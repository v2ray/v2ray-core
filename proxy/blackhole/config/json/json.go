package json

import (
	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
)

type BlackHoleConfig struct {
}

func init() {
	json.RegisterConfigType("blackhole", config.TypeInbound, func() interface{} {
		return new(BlackHoleConfig)
	})
}
