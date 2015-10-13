package vmess

import (
	"net"
	"strings"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
)

// VMessUser is an authenticated user account in VMess configuration.
type VMessUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func (u *VMessUser) ToUser() (user.User, error) {
	id, err := user.NewID(u.Id)
	return user.User{
		Id: id,
	}, err
}

// VMessInboundConfig is
type VMessInboundConfig struct {
	AllowedClients []VMessUser `json:"clients"`
	UDPEnabled     bool        `json:"udp"`
}

type VNextConfig struct {
	Address string      `json:"address"`
	Port    uint16      `json:"port"`
	Users   []VMessUser `json:"users"`
	Network string      `json:"network"`
}

func (config VNextConfig) HasNetwork(network string) bool {
	return strings.Contains(config.Network, network)
}

func (c VNextConfig) ToVNextServer(network string) (*VNextServer, error) {
	users := make([]user.User, 0, len(c.Users))
	for _, user := range c.Users {
		vuser, err := user.ToUser()
		if err != nil {
			log.Error("Failed to convert %v to User.", user)
			return nil, config.BadConfiguration
		}
		users = append(users, vuser)
	}
	ip := net.ParseIP(c.Address)
	if ip == nil {
		log.Error("Unable to parse VNext IP: %s", c.Address)
		return nil, config.BadConfiguration
	}
	address := v2net.IPAddress(ip, c.Port)
	var dest v2net.Destination
	if network == "tcp" {
		dest = v2net.NewTCPDestination(address)
	} else {
		dest = v2net.NewUDPDestination(address)
	}
	return &VNextServer{
		Destination: dest,
		Users:       users,
	}, nil
}

type VMessOutboundConfig struct {
	VNextList []VNextConfig `json:"vnext"`
}

func init() {
	json.RegisterConfigType("vmess", config.TypeInbound, func() interface{} {
		return new(VMessInboundConfig)
	})

	json.RegisterConfigType("vmess", config.TypeOutbound, func() interface{} {
		return new(VMessOutboundConfig)
	})
}
