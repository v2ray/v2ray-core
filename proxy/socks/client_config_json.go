package socks

import (
	"encoding/json"
	"errors"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func (this *ClientConfig) UnmarshalJSON(data []byte) error {
	type ServerConfig struct {
		Address *v2net.AddressJson `json:"address"`
		Port    v2net.Port         `json:"port"`
		Users   []json.RawMessage  `json:"users"`
	}
	type JsonConfig struct {
		Servers []*ServerConfig `json:"servers"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Socks|Client: Failed to parse config: " + err.Error())
	}
	this.Servers = make([]*protocol.ServerSpec, len(jsonConfig.Servers))
	for idx, serverConfig := range jsonConfig.Servers {
		server := protocol.NewServerSpec(v2net.TCPDestination(serverConfig.Address.Address, serverConfig.Port), protocol.AlwaysValid())
		for _, rawUser := range serverConfig.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				return errors.New("Socks|Client: Failed to parse user: " + err.Error())
			}
			account := new(Account)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return errors.New("Socks|Client: Failed to parse socks account: " + err.Error())
			}
			user.Account = account
			server.AddUser(user)
		}
		this.Servers[idx] = server
	}
	return nil
}

func init() {
	internal.RegisterOutboundConfig("socks", func() interface{} { return new(ClientConfig) })
}
