package socks

import (
	"encoding/json"
)

const (
	JsonAuthMethodNoAuth   = "noauth"
	JsonAuthMethodUserPass = "password"
)

type SocksConfig struct {
	AuthMethod string `json:"auth"`
	Username   string `json:"user"`
	Password   string `json:"pass"`
}

func loadConfig(rawConfig []byte) (SocksConfig, error) {
	config := SocksConfig{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}
