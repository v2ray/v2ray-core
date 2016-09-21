// +build json

package internet

import (
	"v2ray.com/core/common/loader"
)

func CreateAuthenticatorConfig(rawConfig []byte) (string, interface{}, error) {
	config, name, err := configCache.Load(rawConfig)
	if err != nil {
		return name, nil, err
	}
	return name, config, nil
}

func init() {
	configCache = loader.NewJSONConfigLoader(loader.ConfigCreatorCache{}, "type", "")
}
