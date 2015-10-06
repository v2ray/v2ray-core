package json

import (
	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
)

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

type SocksConfig struct {
	AuthMethod string `json:"auth"`
	Username   string `json:"user"`
	Password   string `json:"pass"`
	UDPEnabled bool   `json:"udp"`
}

func (config SocksConfig) IsNoAuth() bool {
	return config.AuthMethod == AuthMethodNoAuth
}

func (config SocksConfig) IsPassword() bool {
	return config.AuthMethod == AuthMethodUserPass
}

func init() {
	json.RegisterConfigType("socks", config.TypeInbound, func() interface{} {
		return new(SocksConfig)
	})
}
