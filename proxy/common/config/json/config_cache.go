package json

import (
	"github.com/v2ray/v2ray-core/proxy/common/config"
)

type ConfigObjectCreator func() interface{}

var (
	configCache = make(map[string]ConfigObjectCreator)
)

func getConfigKey(protocol string, cType config.Type) string {
	return protocol + "_" + string(cType)
}

func registerConfigType(protocol string, cType config.Type, creator ConfigObjectCreator) error {
	// TODO: check name
	configCache[getConfigKey(protocol, cType)] = creator
	return nil
}

func RegisterInboundConnectionConfig(protocol string, creator ConfigObjectCreator) error {
	return registerConfigType(protocol, config.TypeInbound, creator)
}

func RegisterOutboundConnectionConfig(protocol string, creator ConfigObjectCreator) error {
	return registerConfigType(protocol, config.TypeOutbound, creator)
}

func CreateConfig(protocol string, cType config.Type) interface{} {
	creator, found := configCache[getConfigKey(protocol, cType)]
	if !found {
		return nil
	}
	return creator()
}
