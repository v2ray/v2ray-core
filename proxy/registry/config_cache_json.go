// +build json

package registry

import (
	"v2ray.com/core/common/loader"
)

func init() {
	inboundConfigCache = loader.NewJSONConfigLoader(loader.ConfigCreatorCache{}, "protocol", "settings")
	outboundConfigCache = loader.NewJSONConfigLoader(loader.ConfigCreatorCache{}, "protocol", "settings")
}
