package conf

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"v2ray.com/core/app/router"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/tools/geoip"

	"github.com/golang/protobuf/proto"
)

type RouterRulesConfig struct {
	RuleList       []json.RawMessage `json:"rules"`
	DomainStrategy string            `json:"domainStrategy"`
}

type RouterConfig struct {
	Settings *RouterRulesConfig `json:"settings"`
}

func (this *RouterConfig) Build() (*router.Config, error) {
	if this.Settings == nil {
		return nil, errors.New("Router settings is not specified.")
	}
	config := new(router.Config)

	settings := this.Settings
	config.DomainStrategy = router.Config_AsIs
	config.Rule = make([]*router.RoutingRule, len(settings.RuleList))
	domainStrategy := strings.ToLower(settings.DomainStrategy)
	if domainStrategy == "alwaysip" {
		config.DomainStrategy = router.Config_UseIp
	} else if domainStrategy == "ipifnonmatch" {
		config.DomainStrategy = router.Config_IpIfNonMatch
	}
	for idx, rawRule := range settings.RuleList {
		rule := ParseRule(rawRule)
		config.Rule[idx] = rule
	}
	return config, nil
}

type RouterRule struct {
	Type        string `json:"type"`
	OutboundTag string `json:"outboundTag"`
}

func parseIP(s string) *router.CIDR {
	var addr, mask string
	i := strings.Index(s, "/")
	if i < 0 {
		addr = s
	} else {
		addr = s[:i]
		mask = s[i+1:]
	}
	ip := v2net.ParseAddress(addr)
	switch ip.Family() {
	case v2net.AddressFamilyIPv4:
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
		return &router.CIDR{
			Ip:     []byte(ip.IP()),
			Prefix: bits,
		}
	case v2net.AddressFamilyIPv6:
		bits := uint32(128)
		if len(mask) > 0 {
			bits64, err := strconv.ParseUint(mask, 10, 32)
			if err != nil {
				return nil
			}
			bits = uint32(bits64)
		}
		if bits > 128 {
			log.Warning("Router: invalid network mask: ", bits)
			return nil
		}
		return &router.CIDR{
			Ip:     []byte(ip.IP()),
			Prefix: bits,
		}
	default:
		log.Warning("Router: unsupported address: ", s)
		return nil
	}
}

func parseFieldRule(msg json.RawMessage) (*router.RoutingRule, error) {
	type RawFieldRule struct {
		RouterRule
		Domain   *StringList  `json:"domain"`
		IP       *StringList  `json:"ip"`
		Port     *PortRange   `json:"port"`
		Network  *NetworkList `json:"network"`
		SourceIP *StringList  `json:"source"`
		User     *StringList  `json:"user"`
	}
	rawFieldRule := new(RawFieldRule)
	err := json.Unmarshal(msg, rawFieldRule)
	if err != nil {
		return nil, err
	}

	rule := new(router.RoutingRule)
	rule.Tag = rawFieldRule.OutboundTag

	if rawFieldRule.Domain != nil {
		for _, domain := range *rawFieldRule.Domain {
			domainRule := new(router.Domain)
			if strings.HasPrefix(domain, "regexp:") {
				domainRule.Type = router.Domain_Regex
				domainRule.Value = domain[7:]
			} else {
				domainRule.Type = router.Domain_Plain
				domainRule.Value = domain
			}
			rule.Domain = append(rule.Domain, domainRule)
		}
	}

	if rawFieldRule.IP != nil {
		for _, ip := range *rawFieldRule.IP {
			ipRule := parseIP(ip)
			if ipRule != nil {
				rule.Cidr = append(rule.Cidr, ipRule)
			}
		}
	}

	if rawFieldRule.Port != nil {
		rule.PortRange = rawFieldRule.Port.Build()
	}

	if rawFieldRule.Network != nil {
		rule.NetworkList = rawFieldRule.Network.Build()
	}

	if rawFieldRule.SourceIP != nil {
		for _, ip := range *rawFieldRule.IP {
			ipRule := parseIP(ip)
			if ipRule != nil {
				rule.SourceCidr = append(rule.SourceCidr, ipRule)
			}
		}
	}

	if rawFieldRule.User != nil {
		for _, s := range *rawFieldRule.User {
			rule.UserEmail = append(rule.UserEmail, s)
		}
	}

	return rule, nil
}

func ParseRule(msg json.RawMessage) *router.RoutingRule {
	rawRule := new(RouterRule)
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

func parseChinaIPRule(data []byte) (*router.RoutingRule, error) {
	rawRule := new(RouterRule)
	err := json.Unmarshal(data, rawRule)
	if err != nil {
		log.Error("Router: Invalid router rule: ", err)
		return nil, err
	}
	var chinaIPs geoip.CountryIPRange
	if err := proto.Unmarshal(geoip.ChinaIPs, &chinaIPs); err != nil {
		return nil, err
	}
	return &router.RoutingRule{
		Tag:  rawRule.OutboundTag,
		Cidr: chinaIPs.Ips,
	}, nil
}

func parseChinaSitesRule(data []byte) (*router.RoutingRule, error) {
	rawRule := new(RouterRule)
	err := json.Unmarshal(data, rawRule)
	if err != nil {
		log.Error("Router: Invalid router rule: ", err)
		return nil, err
	}
	return &router.RoutingRule{
		Tag:    rawRule.OutboundTag,
		Domain: chinaSitesDomains,
	}, nil
}
