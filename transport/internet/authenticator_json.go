// +build json

package internet

import (
	"v2ray.com/core/common/loader"
)

func CreateAuthenticatorConfig(rawConfig []byte) (string, interface{}, error) {
	config, name, err := configLoader.Load(rawConfig)
	if err != nil {
		return name, nil, err
	}
	return name, config, nil
}

var (
	configLoader loader.ConfigLoader
)

func init() {
	configLoader = loader.NewJSONConfigLoader(configCache, "type", "")
}
