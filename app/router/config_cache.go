package router

import (
	"github.com/v2ray/v2ray-core/common"
)

type ConfigObjectCreator func([]byte) (interface{}, error)

var (
	configCache map[string]ConfigObjectCreator
)

func RegisterRouterConfig(strategy string, creator ConfigObjectCreator) error {
	// TODO: check strategy
	configCache[strategy] = creator
	return nil
}

func CreateRouterConfig(strategy string, data []byte) (interface{}, error) {
	creator, found := configCache[strategy]
	if !found {
		return nil, common.ErrObjectNotFound
	}
	return creator(data)
}

func init() {
	configCache = make(map[string]ConfigObjectCreator)
}
