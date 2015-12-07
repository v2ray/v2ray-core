package json

import (
	"encoding/json"
	"net"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
	jsonconfig "github.com/v2ray/v2ray-core/proxy/common/config/json"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	vmessjson "github.com/v2ray/v2ray-core/proxy/vmess/json"
	"github.com/v2ray/v2ray-core/proxy/vmess/outbound"
)

type ConfigTarget struct {
	Address v2net.Address
	Users   []*vmessjson.ConfigUser
}

func (t *ConfigTarget) UnmarshalJSON(data []byte) error {
	type RawConfigTarget struct {
		Address string                  `json:"address"`
		Port    v2net.Port              `json:"port"`
		Users   []*vmessjson.ConfigUser `json:"users"`
	}
	var rawConfig RawConfigTarget
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return err
	}
	if len(rawConfig.Users) == 0 {
		log.Error("0 user configured for VMess outbound.")
		return proxyconfig.BadConfiguration
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

func (this *Outbound) UnmarshalJSON(data []byte) error {
	type RawOutbound struct {
		TargetList []*ConfigTarget `json:"vnext"`
	}
	rawOutbound := &RawOutbound{}
	err := json.Unmarshal(data, rawOutbound)
	if err != nil {
		return err
	}
	if len(rawOutbound.TargetList) == 0 {
		log.Error("0 VMess receiver configured.")
		return proxyconfig.BadConfiguration
	}
	this.TargetList = rawOutbound.TargetList
	return nil
}

func (o *Outbound) Receivers() []*outbound.Receiver {
	targets := make([]*outbound.Receiver, 0, 2*len(o.TargetList))
	for _, rawTarget := range o.TargetList {
		users := make([]vmess.User, 0, len(rawTarget.Users))
		for _, rawUser := range rawTarget.Users {
			users = append(users, rawUser)
		}
		targets = append(targets, &outbound.Receiver{
			Address:  rawTarget.Address,
			Accounts: users,
		})
	}
	return targets
}

func init() {
	jsonconfig.RegisterOutboundConnectionConfig("vmess", func() interface{} {
		return new(Outbound)
	})
}
