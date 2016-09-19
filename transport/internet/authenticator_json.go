// +build json

package internet

import (
	"v2ray.com/core/common/loader"
)

func RegisterAuthenticatorConfig(name string, configCreator loader.ConfigCreator) error {
	return configCache.RegisterCreator(name, configCreator)
}

func CreateAuthenticatorConfig(rawConfig []byte) (string, AuthenticatorConfig, error) {
	config, name, err := configCache.Load(rawConfig)
	if err != nil {
		return name, nil, err
	}
	return name, config, nil
}

var (
	configCache = loader.NewJSONConfigLoader("type", "")
)
