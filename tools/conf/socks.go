package conf

import (
	"encoding/json"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/socks"
)

type SocksAccount struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (v *SocksAccount) Build() *socks.Account {
	return &socks.Account{
		Username: v.Username,
		Password: v.Password,
	}
}

const (
	AuthMethodNoAuth   = "noauth"
	AuthMethodUserPass = "password"
)

type SocksServerConfig struct {
	AuthMethod string          `json:"auth"`
	Accounts   []*SocksAccount `json:"accounts"`
	UDP        bool            `json:"udp"`
	Host       *Address        `json:"ip"`
	Timeout    uint32          `json:"timeout"`
}

func (v *SocksServerConfig) Build() (*loader.TypedSettings, error) {
	config := new(socks.ServerConfig)
	if v.AuthMethod == AuthMethodNoAuth {
		config.AuthType = socks.AuthType_NO_AUTH
	} else if v.AuthMethod == AuthMethodUserPass {
		config.AuthType = socks.AuthType_PASSWORD
	} else {
		return nil, errors.New("Unknown socks auth method: " + v.AuthMethod)
	}

	if len(v.Accounts) > 0 {
		config.Accounts = make(map[string]string, len(v.Accounts))
		for _, account := range v.Accounts {
			config.Accounts[account.Username] = account.Password
		}
	}

	config.UdpEnabled = v.UDP
	if v.Host != nil {
		config.Address = v.Host.Build()
	}

	config.Timeout = v.Timeout
	return loader.NewTypedSettings(config), nil
}

type SocksRemoteConfig struct {
	Address *Address          `json:"address"`
	Port    uint16            `json:"port"`
	Users   []json.RawMessage `json:"users"`
}
type SocksClientConfig struct {
	Servers []*SocksRemoteConfig `json:"servers"`
}

func (v *SocksClientConfig) Build() (*loader.TypedSettings, error) {
	config := new(socks.ClientConfig)
	config.Server = make([]*protocol.ServerEndpoint, len(v.Servers))
	for idx, serverConfig := range v.Servers {
		server := &protocol.ServerEndpoint{
			Address: serverConfig.Address.Build(),
			Port:    uint32(serverConfig.Port),
		}
		for _, rawUser := range serverConfig.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				return nil, errors.New("Socks|Client: Failed to parse user: " + err.Error())
			}
			account := new(SocksAccount)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return nil, errors.New("Socks|Client: Failed to parse socks account: " + err.Error())
			}
			user.Account = loader.NewTypedSettings(account.Build())
			server.User = append(server.User, user)
		}
		config.Server[idx] = server
	}
	return loader.NewTypedSettings(config), nil
}
