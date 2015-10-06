package json

import (
	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
)

type FreedomConfiguration struct {
}

func init() {
	json.RegisterConfigType("freedom", config.TypeOutbound, func() interface{} {
		return &FreedomConfiguration{}
	})
}
