package json

import (
	"encoding/json"
)

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

type SocksConfig struct {
	AuthMethod string `json:"auth"`
	Username   string `json:"user"`
	Password   string `json:"pass"`
}

func (config SocksConfig) IsNoAuth() bool {
	return config.AuthMethod == AuthMethodNoAuth
}

func (config SocksConfig) IsPassword() bool {
	return config.AuthMethod == AuthMethodUserPass
}

func Load(rawConfig []byte) (SocksConfig, error) {
	config := SocksConfig{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}
