package json

import (
	"encoding/json"
	"errors"
	"net"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netjson "github.com/v2ray/v2ray-core/common/net/json"
)

type FieldRule struct {
	Rule
	Domain  string
	IP      *net.IPNet
	Port    v2net.PortRange
	Network v2net.NetworkList
}

func (this *FieldRule) Apply(dest v2net.Destination) bool {
	address := dest.Address()
	if len(this.Domain) > 0 {
		if !address.IsDomain() || !strings.Contains(address.Domain(), this.Domain) {
			return false
		}
	}

	if this.IP != nil {
		if !(address.IsIPv4() || address.IsIPv6()) || !this.IP.Contains(address.IP()) {
			return false
		}
	}

	if this.Port != nil {
		port := address.Port()
		if port < this.Port.From() || port > this.Port.To() {
			return false
		}
	}

	if this.Network != nil {
		if !this.Network.HasNetwork(v2net.Network(dest.Network())) {
			return false
		}
	}

	return true
}

func (this *FieldRule) UnmarshalJSON(data []byte) error {
	type RawFieldRule struct {
		Rule
		Domain  string                 `json:"domain"`
		IP      string                 `json:"ip"`
		Port    *v2netjson.PortRange   `json:"port"`
		Network *v2netjson.NetworkList `json:"network"`
	}
	rawFieldRule := RawFieldRule{}
	err := json.Unmarshal(data, &rawFieldRule)
	if err != nil {
		return err
	}
	this.Type = rawFieldRule.Type
	this.OutboundTag = rawFieldRule.OutboundTag

	hasField := false
	if len(rawFieldRule.Domain) > 0 {
		this.Domain = rawFieldRule.Domain
		hasField = true
	}

	if len(rawFieldRule.IP) > 0 {
		_, ipNet, err := net.ParseCIDR(rawFieldRule.IP)
		if err != nil {
			return errors.New("Invalid IP range in router rule: " + err.Error())
		}
		this.IP = ipNet
		hasField = true
	}
	if rawFieldRule.Port != nil {
		this.Port = rawFieldRule.Port
		hasField = true
	}
	if rawFieldRule.Network != nil {
		this.Network = rawFieldRule.Network
		hasField = true
	}
	if !hasField {
		return errors.New("This rule has no effective fields.")
	}
	return nil
}
