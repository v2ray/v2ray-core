package json

import (
	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
)

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

type SocksAccount struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

type SocksConfig struct {
	AuthMethod string         `json:"auth"`
	Accounts   []SocksAccount `json:"accounts"`
	UDPEnabled bool           `json:"udp"`

	accountMap map[string]string
}

func (config *SocksConfig) Initialize() {
	config.accountMap = make(map[string]string)
	for _, account := range config.Accounts {
		config.accountMap[account.Username] = account.Password
	}
}

func (config *SocksConfig) IsNoAuth() bool {
	return config.AuthMethod == AuthMethodNoAuth
}

func (config *SocksConfig) IsPassword() bool {
	return config.AuthMethod == AuthMethodUserPass
}

func (config *SocksConfig) HasAccount(user, pass string) bool {
	if actualPass, found := config.accountMap[user]; found {
		return actualPass == pass
	}
	return false
}

func init() {
	json.RegisterConfigType("socks", config.TypeInbound, func() interface{} {
		return new(SocksConfig)
	})
}
