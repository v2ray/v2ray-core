// +build json

package registry

import (
	"v2ray.com/core/common/loader"
)

var (
	inboundConfigCache  loader.ConfigLoader
	outboundConfigCache loader.ConfigLoader
)

func CreateInboundConfig(protocol string, data []byte) (interface{}, error) {
	return inboundConfigCache.LoadWithID(data, protocol)
}

func CreateOutboundConfig(protocol string, data []byte) (interface{}, error) {
	return outboundConfigCache.LoadWithID(data, protocol)
}

func init() {
	inboundConfigCache = loader.NewJSONConfigLoader(inboundConfigCreatorCache, "protocol", "settings")
	outboundConfigCache = loader.NewJSONConfigLoader(outboundConfigCreatorCache, "protocol", "settings")
}
