// +build json

package socks

import (
	"encoding/json"
	"errors"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/registry"
)

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type SocksConfig struct {
		AuthMethod string             `json:"auth"`
		Accounts   []*Account         `json:"accounts"`
		UDP        bool               `json:"udp"`
		Host       *v2net.AddressJson `json:"ip"`
		Timeout    int                `json:"timeout"`
	}

	rawConfig := new(SocksConfig)
	if err := json.Unmarshal(data, rawConfig); err != nil {
		return errors.New("Socks: Failed to parse config: " + err.Error())
	}
	if rawConfig.AuthMethod == AuthMethodNoAuth {
		this.AuthType = AuthTypeNoAuth
	} else if rawConfig.AuthMethod == AuthMethodUserPass {
		this.AuthType = AuthTypePassword
	} else {
		log.Error("Socks: Unknown auth method: ", rawConfig.AuthMethod)
		return common.ErrBadConfiguration
	}

	if len(rawConfig.Accounts) > 0 {
		this.Accounts = make(map[string]string, len(rawConfig.Accounts))
		for _, account := range rawConfig.Accounts {
			this.Accounts[account.Username] = account.Password
		}
	}

	this.UDPEnabled = rawConfig.UDP
	if rawConfig.Host != nil {
		this.Address = rawConfig.Host.Address
	} else {
		this.Address = v2net.LocalHostIP
	}

	if rawConfig.Timeout >= 0 {
		this.Timeout = rawConfig.Timeout
	}
	return nil
}

func init() {
	registry.RegisterInboundConfig("socks", func() interface{} { return new(Config) })
}
