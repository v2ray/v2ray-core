package vmess

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/log"
	v2net "github.com/v2ray/v2ray-core/net"
)

type VMessUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func (u *VMessUser) ToVUser() (core.VUser, error) {
	id, err := core.UUIDToVID(u.Id)
	return core.VUser{id}, err
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
	users := make([]core.VUser, 0, len(config.Users))
	for _, user := range config.Users {
		vuser, err := user.ToVUser()
		if err != nil {
			panic(log.Error("Failed to convert %v to VUser.", user))
		}
		users = append(users, vuser)
	}
	return VNextServer{
		v2net.DomainAddress(config.Address, config.Port),
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
