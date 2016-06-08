// +build json

package blackhole

import (
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

func init() {
	config.RegisterOutboundConfig("blackhole",
		func(data []byte) (interface{}, error) {
			return new(Config), nil
		})
}
