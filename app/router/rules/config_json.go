// +build json

package rules

import (
	"encoding/json"
	"strconv"
	"strings"

	router "v2ray.com/core/app/router"
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

func init() {
	router.RegisterRouterConfig("rules", func(data []byte) (interface{}, error) {
		type JsonConfig struct {
			RuleList       []json.RawMessage `json:"rules"`
			DomainStrategy string            `json:"domainStrategy"`
		}
		jsonConfig := new(JsonConfig)
		if err := json.Unmarshal(data, jsonConfig); err != nil {
			return nil, err
		}
		config := &Config{
			Rule:           make([]*RoutingRule, len(jsonConfig.RuleList)),
			DomainStrategy: Config_AsIs,
		}
		domainStrategy := strings.ToLower(jsonConfig.DomainStrategy)
		if domainStrategy == "alwaysip" {
			config.DomainStrategy = Config_UseIp
		} else if domainStrategy == "ipifnonmatch" {
			config.DomainStrategy = Config_IpIfNonMatch
		}
		for idx, rawRule := range jsonConfig.RuleList {
			rule := ParseRule(rawRule)
			config.Rule[idx] = rule
		}
		return config, nil
	})
}
