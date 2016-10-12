// +build json

package router

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"v2ray.com/core/common/collect"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
)

type JsonRule struct {
	Type        string `json:"type"`
	OutboundTag string `json:"outboundTag"`
}

func parseIP(s string) *IP {
	var addr, mask string
	i := strings.Index(s, "/")
	if i < 0 {
		addr = s
	} else {
		addr = s[:i]
		mask = s[i+1:]
	}
	ip := v2net.ParseAddress(addr)
	if !ip.Family().Either(v2net.AddressFamilyIPv4, v2net.AddressFamilyIPv6) {
		return nil
	}
	bits := uint32(32)
	if len(mask) > 0 {
		bits64, err := strconv.ParseUint(mask, 10, 32)
		if err != nil {
			return nil
		}
		bits = uint32(bits64)
	}
	if bits > 32 {
		log.Warning("Router: invalid network mask: ", bits)
		return nil
	}
	return &IP{
		Ip:             []byte(ip.IP()),
		UnmatchingBits: 32 - bits,
	}
}

func parseFieldRule(msg json.RawMessage) (*RoutingRule, error) {
	type RawFieldRule struct {
		JsonRule
		Domain  *collect.StringList `json:"domain"`
		IP      *collect.StringList `json:"ip"`
		Port    *v2net.PortRange    `json:"port"`
		Network *v2net.NetworkList  `json:"network"`
	}
	rawFieldRule := new(RawFieldRule)
	err := json.Unmarshal(msg, rawFieldRule)
	if err != nil {
		return nil, err
	}

	rule := new(RoutingRule)
	rule.Tag = rawFieldRule.OutboundTag

	if rawFieldRule.Domain != nil {
		for _, domain := range *rawFieldRule.Domain {
			domainRule := new(Domain)
			if strings.HasPrefix(domain, "regexp:") {
				domainRule.Type = Domain_Regex
				domainRule.Value = domain[7:]
			} else {
				domainRule.Type = Domain_Plain
				domainRule.Value = domain
			}
			rule.Domain = append(rule.Domain, domainRule)
		}
	}

	if rawFieldRule.IP != nil {
		for _, ip := range *rawFieldRule.IP {
			ipRule := parseIP(ip)
			if ipRule != nil {
				rule.Ip = append(rule.Ip, ipRule)
			}
		}
	}

	if rawFieldRule.Port != nil {
		rule.PortRange = rawFieldRule.Port
	}

	if rawFieldRule.Network != nil {
		rule.NetworkList = rawFieldRule.Network
	}

	return rule, nil
}

func ParseRule(msg json.RawMessage) *RoutingRule {
	rawRule := new(JsonRule)
	err := json.Unmarshal(msg, rawRule)
	if err != nil {
		log.Error("Router: Invalid router rule: ", err)
		return nil
	}
	if rawRule.Type == "field" {

		fieldrule, err := parseFieldRule(msg)
		if err != nil {
			log.Error("Invalid field rule: ", err)
			return nil
		}
		return fieldrule
	}
	if rawRule.Type == "chinaip" {
		chinaiprule, err := parseChinaIPRule(msg)
		if err != nil {
			log.Error("Router: Invalid chinaip rule: ", err)
			return nil
		}
		return chinaiprule
	}
	if rawRule.Type == "chinasites" {
		chinasitesrule, err := parseChinaSitesRule(msg)
		if err != nil {
			log.Error("Invalid chinasites rule: ", err)
			return nil
		}
		return chinasitesrule
	}
	log.Error("Unknown router rule type: ", rawRule.Type)
	return nil
}

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonRulesConfig struct {
		RuleList       []json.RawMessage `json:"rules"`
		DomainStrategy string            `json:"domainStrategy"`
	}
	type JsonConfig struct {
		Settings *JsonRulesConfig `json:"settings"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	if jsonConfig.Settings == nil {
		return errors.New("Router settings is not specified.")
	}
	settings := jsonConfig.Settings
	this.DomainStrategy = Config_AsIs
	this.Rule = make([]*RoutingRule, len(settings.RuleList))
	domainStrategy := strings.ToLower(settings.DomainStrategy)
	if domainStrategy == "alwaysip" {
		this.DomainStrategy = Config_UseIp
	} else if domainStrategy == "ipifnonmatch" {
		this.DomainStrategy = Config_IpIfNonMatch
	}
	for idx, rawRule := range settings.RuleList {
		rule := ParseRule(rawRule)
		this.Rule[idx] = rule
	}
	return nil
}
