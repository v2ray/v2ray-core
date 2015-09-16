package vmess

import (
	"encoding/json"
	"net"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/log"
	v2net "github.com/v2ray/v2ray-core/net"
)

type VMessUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func (u *VMessUser) ToUser() (core.User, error) {
	id, err := core.NewID(u.Id)
	return core.User{id}, err
}

type VMessInboundConfig struct {
	AllowedClients []VMessUser `json:"clients"`
}

func loadInboundConfig(rawConfig []byte) (VMessInboundConfig, error) {
	config := VMessInboundConfig{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}

type VNextConfig struct {
	Address string      `json:"address"`
	Port    uint16      `json:"port"`
	Users   []VMessUser `json:"users"`
}

func (config VNextConfig) ToVNextServer() VNextServer {
	users := make([]core.User, 0, len(config.Users))
	for _, user := range config.Users {
		vuser, err := user.ToUser()
		if err != nil {
			panic(log.Error("Failed to convert %v to User.", user))
		}
		users = append(users, vuser)
	}
	ip := net.ParseIP(config.Address)
	if ip == nil {
		panic(log.Error("Unable to parse VNext IP: %s", config.Address))
	}
	return VNextServer{
		v2net.IPAddress(ip, config.Port),
		users}
}

type VMessOutboundConfig struct {
	VNextList []VNextConfig `json:"vnext"`
}

func loadOutboundConfig(rawConfig []byte) (VMessOutboundConfig, error) {
	config := VMessOutboundConfig{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}
