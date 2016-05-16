// +build json

package rules

import (
	"encoding/json"
	"errors"
	"strings"

	router "github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
)

type JsonRule struct {
	Type        string `json:"type"`
	OutboundTag string `json:"outboundTag"`
}

func parseFieldRule(msg json.RawMessage) (*Rule, error) {
	type RawFieldRule struct {
		JsonRule
		Domain  *serial.StringLiteralList `json:"domain"`
		IP      *serial.StringLiteralList `json:"ip"`
		Port    *v2net.PortRange          `json:"port"`
		Network *v2net.NetworkList        `json:"network"`
	}
	rawFieldRule := new(RawFieldRule)
	err := json.Unmarshal(msg, rawFieldRule)
	if err != nil {
		return nil, err
	}
	conds := NewConditionChan()

	if rawFieldRule.Domain != nil && rawFieldRule.Domain.Len() > 0 {
		anyCond := NewAnyCondition()
		for _, rawDomain := range *(rawFieldRule.Domain) {
			var matcher Condition
			if strings.HasPrefix(rawDomain.String(), "regexp:") {
				rawMatcher, err := NewRegexpDomainMatcher(rawDomain.String()[7:])
				if err != nil {
					return nil, err
				}
				matcher = rawMatcher
			} else {
				matcher = NewPlainDomainMatcher(rawDomain.String())
			}
			anyCond.Add(matcher)
		}
		conds.Add(anyCond)
	}

	if rawFieldRule.IP != nil && rawFieldRule.IP.Len() > 0 {
		anyCond := NewAnyCondition()
		for _, ipStr := range *(rawFieldRule.IP) {
			cidrMatcher, err := NewCIDRMatcher(ipStr.String())
			if err != nil {
				log.Error("Router: Invalid IP range in router rule: ", err)
				return nil, err
			}
			anyCond.Add(cidrMatcher)
		}
		conds.Add(anyCond)
	}
	if rawFieldRule.Port != nil {
		conds.Add(NewPortMatcher(*rawFieldRule.Port))
	}
	if rawFieldRule.Network != nil {
		conds.Add(NewNetworkMatcher(rawFieldRule.Network))
	}
	if conds.Len() == 0 {
		return nil, errors.New("Router: This rule has no effective fields.")
	}
	return &Rule{
		Tag:       rawFieldRule.OutboundTag,
		Condition: conds,
	}, nil
}

func ParseRule(msg json.RawMessage) *Rule {
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
		config := &RouterRuleConfig{
			Rules:          make([]*Rule, len(jsonConfig.RuleList)),
			DomainStrategy: DomainAsIs,
		}
		domainStrategy := serial.StringLiteral(jsonConfig.DomainStrategy).ToLower()
		if domainStrategy.String() == "alwaysip" {
			config.DomainStrategy = AlwaysUseIP
		} else if domainStrategy.String() == "ipifnonmatch" {
			config.DomainStrategy = UseIPIfNonMatch
		}
		for idx, rawRule := range jsonConfig.RuleList {
			rule := ParseRule(rawRule)
			config.Rules[idx] = rule
		}
		return config, nil
	})
}
