package json

import (
	"github.com/v2ray/v2ray-core/config"
)

type ConfigObjectCreator func() interface{}

var (
	configCache = make(map[string]ConfigObjectCreator)
)

func getConfigKey(protocol string, cType config.Type) string {
	return protocol + "_" + string(cType)
}

func RegisterConfigType(protocol string, cType config.Type, creator ConfigObjectCreator) {
	// TODO: check name
	configCache[getConfigKey(protocol, cType)] = creator
}
