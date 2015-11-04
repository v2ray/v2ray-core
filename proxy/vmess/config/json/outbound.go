package json

import (
	"encoding/json"
	"net"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
	jsonconfig "github.com/v2ray/v2ray-core/proxy/common/config/json"
	vmessconfig "github.com/v2ray/v2ray-core/proxy/vmess/config"
)

type RawConfigTarget struct {
	Address string        `json:"address"`
	Port    uint16        `json:"port"`
	Users   []*ConfigUser `json:"users"`
}

type ConfigTarget struct {
	Address v2net.Address
	Users   []*ConfigUser
}

func (t *ConfigTarget) UnmarshalJSON(data []byte) error {
	var rawConfig RawConfigTarget
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return err
	}
	t.Users = rawConfig.Users
	ip := net.ParseIP(rawConfig.Address)
	if ip == nil {
		log.Error("Unable to parse IP: %s", rawConfig.Address)
		return proxyconfig.BadConfiguration
	}
	t.Address = v2net.IPAddress(ip, rawConfig.Port)
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
		targets = append(targets, &vmessconfig.OutboundTarget{
			Destination: v2net.NewTCPDestination(rawTarget.Address),
			Accounts:    users,
		})
	}
	return targets
}

func init() {
	jsonconfig.RegisterOutboundConnectionConfig("vmess", func() interface{} {
		return new(Outbound)
	})
}
