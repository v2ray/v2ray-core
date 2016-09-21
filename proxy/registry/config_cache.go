package registry

import "v2ray.com/core/common/loader"

var (
	inboundConfigCreatorCache = loader.ConfigCreatorCache{}
	inboundConfigCache        loader.ConfigLoader

	outboundConfigCreatorCache = loader.ConfigCreatorCache{}
	outboundConfigCache        loader.ConfigLoader
)

func RegisterInboundConfig(protocol string, creator loader.ConfigCreator) error {
	return inboundConfigCreatorCache.RegisterCreator(protocol, creator)
}

func RegisterOutboundConfig(protocol string, creator loader.ConfigCreator) error {
	return outboundConfigCreatorCache.RegisterCreator(protocol, creator)
}

func CreateInboundConfig(protocol string, data []byte) (interface{}, error) {
	return inboundConfigCache.LoadWithID(data, protocol)
}

func CreateOutboundConfig(protocol string, data []byte) (interface{}, error) {
	return outboundConfigCache.LoadWithID(data, protocol)
}
