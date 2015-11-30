package json

import (
	"encoding/json"
	"errors"
	"net"
	"strings"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2netjson "github.com/v2ray/v2ray-core/common/net/json"
)

type StringList []string

func NewStringList(str ...string) *StringList {
	list := StringList(str)
	return &list
}

func (this *StringList) UnmarshalJSON(data []byte) error {
	var strList []string
	err := json.Unmarshal(data, &strList)
	if err == nil {
		*this = make([]string, len(strList))
		copy(*this, strList)
		return nil
	}

	var str string
	err = json.Unmarshal(data, &str)
	if err == nil {
		*this = make([]string, 0, 1)
		*this = append(*this, str)
		return nil
	}

	return errors.New("Failed to unmarshal string list: " + string(data))
}

func (this *StringList) Len() int {
	return len([]string(*this))
}

type FieldRule struct {
	Rule
	Domain  *StringList
	IP      []*net.IPNet
	Port    v2net.PortRange
	Network v2net.NetworkList
}

func (this *FieldRule) Apply(dest v2net.Destination) bool {
	address := dest.Address()
	if this.Domain != nil && this.Domain.Len() > 0 {
		if !address.IsDomain() {
			return false
		}
		foundMatch := false
		for _, domain := range *this.Domain {
			if strings.Contains(address.Domain(), domain) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
			return false
		}
	}

	if this.IP != nil && len(this.IP) > 0 {
		if !(address.IsIPv4() || address.IsIPv6()) {
			return false
		}
		foundMatch := false
		for _, ipnet := range this.IP {
			if ipnet.Contains(address.IP()) {
				foundMatch = true
				break
			}
		}
		if !foundMatch {
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
		Domain  *StringList            `json:"domain"`
		IP      *StringList            `json:"ip"`
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
	if rawFieldRule.Domain != nil && rawFieldRule.Domain.Len() > 0 {
		this.Domain = rawFieldRule.Domain
		hasField = true
	}

	if rawFieldRule.IP != nil && rawFieldRule.IP.Len() > 0 {
		this.IP = make([]*net.IPNet, 0, rawFieldRule.IP.Len())
		for _, ipStr := range *(rawFieldRule.IP) {
			_, ipNet, err := net.ParseCIDR(ipStr)
			if err != nil {
				return errors.New("Invalid IP range in router rule: " + err.Error())
			}
			this.IP = append(this.IP, ipNet)
		}
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
