package json

import (
	"encoding/json"
	"net"
	"strings"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/config"
	jsonconfig "github.com/v2ray/v2ray-core/config/json"
	vmessconfig "github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type RawConfigTarget struct {
	Address string        `json:"address"`
	Port    uint16        `json:"port"`
	Users   []*ConfigUser `json:"users"`
	Network string        `json:"network"`
}

func (config RawConfigTarget) HasNetwork(network string) bool {
	return strings.Contains(config.Network, network)
}

type ConfigTarget struct {
	Address    v2net.Address
	Users      []*ConfigUser
	TCPEnabled bool
	UDPEnabled bool
}

func (t *ConfigTarget) UnmarshalJSON(data []byte) error {
	var rawConfig RawConfigTarget
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return err
	}
	ip := net.ParseIP(rawConfig.Address)
	if ip == nil {
		log.Error("Unable to parse IP: %s", rawConfig.Address)
		return config.BadConfiguration
	}
	t.Address = v2net.IPAddress(ip, rawConfig.Port)
	if rawConfig.HasNetwork("tcp") {
		t.TCPEnabled = true
	}
	if rawConfig.HasNetwork("udp") {
		t.UDPEnabled = true
	}
	return nil
}

type Outbound struct {
	TargetList []*ConfigTarget `json:"vnext"`
}

func (o *Outbound) Targets() []*vmessconfig.OutboundTarget {
	targets := make([]*vmessconfig.OutboundTarget, 0, 2*len(o.TargetList))
	for _, rawTarget := range o.TargetList {
		users := make([]vmessconfig.User, 0, len(rawTarget.Users))
		for _, rawUser := range rawTarget.Users {
			users = append(users, rawUser)
		}
		if rawTarget.TCPEnabled {
			targets = append(targets, &vmessconfig.OutboundTarget{
				Destination: v2net.NewTCPDestination(rawTarget.Address),
				Accounts:    users,
			})
		}
		if rawTarget.UDPEnabled {
			targets = append(targets, &vmessconfig.OutboundTarget{
				Destination: v2net.NewUDPDestination(rawTarget.Address),
				Accounts:    users,
			})
		}
	}
	return targets
}

func init() {
	jsonconfig.RegisterConfigType("vmess", config.TypeOutbound, func() interface{} {
		return new(Outbound)
	})
}
