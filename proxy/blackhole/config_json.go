// +build json

package blackhole

import (
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	return nil
}

func init() {
	internal.RegisterOutboundConfig("blackhole", func() interface{} { return new(Config) })
}
