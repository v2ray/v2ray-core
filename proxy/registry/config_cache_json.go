// +build json

package registry

import (
	"v2ray.com/core/common/loader"
)

func init() {
	inboundConfigCache = loader.NewJSONConfigLoader("protocol", "settings")
	outboundConfigCache = loader.NewJSONConfigLoader("protocol", "settings")
}
