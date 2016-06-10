package internal

import "github.com/v2ray/v2ray-core/common/loader"

var (
	inboundConfigCache  loader.ConfigLoader
	outboundConfigCache loader.ConfigLoader
)

func RegisterInboundConfig(protocol string, creator loader.ConfigCreator) error {
	return inboundConfigCache.RegisterCreator(protocol, creator)
}

func RegisterOutboundConfig(protocol string, creator loader.ConfigCreator) error {
	return outboundConfigCache.RegisterCreator(protocol, creator)
}

func CreateInboundConfig(protocol string, data []byte) (interface{}, error) {
	return inboundConfigCache.LoadWithID(data, protocol)
}

func CreateOutboundConfig(protocol string, data []byte) (interface{}, error) {
	return outboundConfigCache.LoadWithID(data, protocol)
}
