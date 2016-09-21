// +build json

package registry

import (
	"v2ray.com/core/common/loader"
)

func init() {
	inboundConfigCache = loader.NewJSONConfigLoader(inboundConfigCreatorCache, "protocol", "settings")
	outboundConfigCache = loader.NewJSONConfigLoader(outboundConfigCreatorCache, "protocol", "settings")
}
