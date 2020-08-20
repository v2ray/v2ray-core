package conf

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/http"
)

type HttpAccount struct {
	Username string `json:"user"`
	Password string `json:"pass"`
}

func (v *HttpAccount) Build() *http.Account {
	return &http.Account{
		Username: v.Username,
		Password: v.Password,
	}
}

type HttpServerConfig struct {
	Timeout     uint32         `json:"timeout"`
	Accounts    []*HttpAccount `json:"accounts"`
	Transparent bool           `json:"allowTransparent"`
	UserLevel   uint32         `json:"userLevel"`
}

func (c *HttpServerConfig) Build() (proto.Message, error) {
	config := &http.ServerConfig{
		Timeout:          c.Timeout,
		AllowTransparent: c.Transparent,
		UserLevel:        c.UserLevel,
	}

	if len(c.Accounts) > 0 {
		config.Accounts = make(map[string]string)
		for _, account := range c.Accounts {
			config.Accounts[account.Username] = account.Password
		}
	}

	return config, nil
}

type HttpRemoteConfig struct {
	Address *Address          `json:"address"`
	Port    uint16            `json:"port"`
	Users   []json.RawMessage `json:"users"`
}
type HttpClientConfig struct {
	Servers []*HttpRemoteConfig `json:"servers"`
}

func (v *HttpClientConfig) Build() (proto.Message, error) {
	config := new(http.ClientConfig)
	config.Server = make([]*protocol.ServerEndpoint, len(v.Servers))
	for idx, serverConfig := range v.Servers {
		server := &protocol.ServerEndpoint{
			Address: serverConfig.Address.Build(),
			Port:    uint32(serverConfig.Port),
		}
		for _, rawUser := range serverConfig.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				return nil, newError("failed to parse HTTP user").Base(err).AtError()
			}
			account := new(HttpAccount)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return nil, newError("failed to parse HTTP account").Base(err).AtError()
			}
			user.Account = serial.ToTypedMessage(account.Build())
			server.User = append(server.User, user)
		}
		config.Server[idx] = server
	}
	return config, nil
}
