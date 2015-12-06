package json

import (
	"github.com/v2ray/v2ray-core/proxy/common/config/json"
)

type BlackHoleConfig struct {
}

func init() {
	json.RegisterOutboundConnectionConfig("blackhole", func() interface{} {
		return new(BlackHoleConfig)
	})
}
