// +build json

package socks

import (
	"encoding/json"
	"errors"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/registry"
)

func (this *Account) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Username string `json:"user"`
		Password string `json:"pass"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Socks: Failed to parse account: " + err.Error())
	}
	this.Username = jsonConfig.Username
	this.Password = jsonConfig.Password
	return nil
}

func (this *ClientConfig) UnmarshalJSON(data []byte) error {
	type ServerConfig struct {
		Address *v2net.AddressPB  `json:"address"`
		Port    v2net.Port        `json:"port"`
		Users   []json.RawMessage `json:"users"`
	}
	type JsonConfig struct {
		Servers []*ServerConfig `json:"servers"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return errors.New("Socks|Client: Failed to parse config: " + err.Error())
	}
	this.Server = make([]*protocol.ServerSpecPB, len(jsonConfig.Servers))
	for idx, serverConfig := range jsonConfig.Servers {
		server := &protocol.ServerSpecPB{
			Address: serverConfig.Address,
			Port:    uint32(serverConfig.Port),
		}
		for _, rawUser := range serverConfig.Users {
			user := new(protocol.User)
			if err := json.Unmarshal(rawUser, user); err != nil {
				return errors.New("Socks|Client: Failed to parse user: " + err.Error())
			}
			account := new(Account)
			if err := json.Unmarshal(rawUser, account); err != nil {
				return errors.New("Socks|Client: Failed to parse socks account: " + err.Error())
			}
			anyAccount, err := account.AsAny()
			if err != nil {
				return err
			}
			user.Account = anyAccount
			server.Users = append(server.Users, user)
		}
		this.Server[idx] = server
	}
	return nil
}

func init() {
	registry.RegisterOutboundConfig("socks", func() interface{} { return new(ClientConfig) })
}
