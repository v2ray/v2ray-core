// +build json

package socks

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netjson "github.com/v2ray/v2ray-core/common/net/json"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/proxy/internal/config"
)

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

func init() {
	config.RegisterInboundConnectionConfig("socks",
		func(data []byte) (interface{}, error) {
			type SocksAccount struct {
				Username string `json:"user"`
				Password string `json:"pass"`
			}

			type SocksConfig struct {
				AuthMethod string          `json:"auth"`
				Accounts   []*SocksAccount `json:"accounts"`
				UDP        bool            `json:"udp"`
				Host       *v2netjson.Host `json:"ip"`
			}

			rawConfig := new(SocksConfig)
			if err := json.Unmarshal(data, rawConfig); err != nil {
				return nil, err
			}
			socksConfig := new(Config)
			if rawConfig.AuthMethod == AuthMethodNoAuth {
				socksConfig.AuthType = AuthTypeNoAuth
			} else if rawConfig.AuthMethod == AuthMethodUserPass {
				socksConfig.AuthType = AuthTypePassword
			} else {
				log.Error("Socks: Unknown auth method: %s", rawConfig.AuthMethod)
				return nil, internal.ErrorBadConfiguration
			}

			if len(rawConfig.Accounts) > 0 {
				socksConfig.Accounts = make(map[string]string, len(rawConfig.Accounts))
				for _, account := range rawConfig.Accounts {
					socksConfig.Accounts[account.Username] = account.Password
				}
			}

			socksConfig.UDPEnabled = rawConfig.UDP
			if rawConfig.Host != nil {
				socksConfig.Address = rawConfig.Host.Address()
			} else {
				socksConfig.Address = v2net.IPAddress([]byte{127, 0, 0, 1})
			}
			return socksConfig, nil
		})
}
