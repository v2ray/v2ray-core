package json

import (
	"github.com/v2ray/v2ray-core/proxy/common/config/json"
)

type FreedomConfiguration struct {
}

func init() {
	json.RegisterOutboundConnectionConfig("freedom", func() interface{} {
		return &FreedomConfiguration{}
	})
}
