package conf

import (
	"encoding/json"
	"strconv"
	"strings"

	"v2ray.com/core/app/log"
	"v2ray.com/core/app/router"
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

func (c *RouterConfig) Build() (*router.Config, error) {
	if c.Settings == nil {
		return nil, newError("Router settings is not specified.")
	}
	config := new(router.Config)

	settings := c.Settings
	config.DomainStrategy = router.Config_AsIs
	config.Rule = make([]*router.RoutingRule, len(settings.RuleList))
	domainStrategy := strings.ToLower(settings.DomainStrategy)
	if domainStrategy == "alwaysip" {
		config.DomainStrategy = router.Config_UseIp
	} else if domainStrategy == "ipifnonmatch" {
		config.DomainStrategy = router.Config_IpIfNonMatch
	}
	for idx, rawRule := range settings.RuleList {
		rule, err := ParseRule(rawRule)
		if err != nil {
			return nil, err
		}
		config.Rule[idx] = rule
	}
	return config, nil
}

type RouterRule struct {
	Type        string `json:"type"`
	OutboundTag string `json:"outboundTag"`
}

func parseIP(s string) (*router.CIDR, error) {
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
				return nil, newError("invalid network mask for router: ", mask).Base(err)
			}
			bits = uint32(bits64)
		}
		if bits > 32 {
			return nil, newError("invalid network mask for router: ", bits)
		}
		return &router.CIDR{
			Ip:     []byte(ip.IP()),
			Prefix: bits,
		}, nil
	case v2net.AddressFamilyIPv6:
		bits := uint32(128)
		if len(mask) > 0 {
			bits64, err := strconv.ParseUint(mask, 10, 32)
			if err != nil {
				return nil, newError("invalid network mask for router: ", mask).Base(err)
			}
			bits = uint32(bits64)
		}
		if bits > 128 {
			return nil, newError("invalid network mask for router: ", bits)
		}
		return &router.CIDR{
			Ip:     []byte(ip.IP()),
			Prefix: bits,
		}, nil
	default:
		return nil, newError("unsupported address for router: ", s)
	}
}

func parseFieldRule(msg json.RawMessage) (*router.RoutingRule, error) {
	type RawFieldRule struct {
		RouterRule
		Domain     *StringList  `json:"domain"`
		IP         *StringList  `json:"ip"`
		Port       *PortRange   `json:"port"`
		Network    *NetworkList `json:"network"`
		SourceIP   *StringList  `json:"source"`
		User       *StringList  `json:"user"`
		InboundTag *StringList  `json:"inboundTag"`
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
			switch {
			case strings.HasPrefix(domain, "regexp:"):
				domainRule.Type = router.Domain_Regex
				domainRule.Value = domain[7:]
			case strings.HasPrefix(domain, "domain:"):
				domainRule.Type = router.Domain_Domain
				domainRule.Value = domain[7:]
			default:
				domainRule.Type = router.Domain_Plain
				domainRule.Value = domain
			}
			rule.Domain = append(rule.Domain, domainRule)
		}
	}

	if rawFieldRule.IP != nil {
		for _, ip := range *rawFieldRule.IP {
			ipRule, err := parseIP(ip)
			if err != nil {
				return nil, newError("invalid IP: ", ip).Base(err)
			}
			rule.Cidr = append(rule.Cidr, ipRule)
		}
	}

	if rawFieldRule.Port != nil {
		rule.PortRange = rawFieldRule.Port.Build()
	}

	if rawFieldRule.Network != nil {
		rule.NetworkList = rawFieldRule.Network.Build()
	}

	if rawFieldRule.SourceIP != nil {
		for _, ip := range *rawFieldRule.SourceIP {
			ipRule, err := parseIP(ip)
			if err != nil {
				return nil, newError("invalid IP: ", ip).Base(err)
			}
			rule.SourceCidr = append(rule.SourceCidr, ipRule)
		}
	}

	if rawFieldRule.User != nil {
		for _, s := range *rawFieldRule.User {
			rule.UserEmail = append(rule.UserEmail, s)
		}
	}

	if rawFieldRule.InboundTag != nil {
		for _, s := range *rawFieldRule.InboundTag {
			rule.InboundTag = append(rule.InboundTag, s)
		}
	}

	return rule, nil
}

func ParseRule(msg json.RawMessage) (*router.RoutingRule, error) {
	rawRule := new(RouterRule)
	err := json.Unmarshal(msg, rawRule)
	if err != nil {
		return nil, newError("invalid router rule").Base(err)
	}
	if rawRule.Type == "field" {
		fieldrule, err := parseFieldRule(msg)
		if err != nil {
			return nil, newError("invalid field rule").Base(err)
		}
		return fieldrule, nil
	}
	if rawRule.Type == "chinaip" {
		chinaiprule, err := parseChinaIPRule(msg)
		if err != nil {
			return nil, newError("invalid chinaip rule").Base(err)
		}
		return chinaiprule, nil
	}
	if rawRule.Type == "chinasites" {
		chinasitesrule, err := parseChinaSitesRule(msg)
		if err != nil {
			return nil, newError("invalid chinasites rule").Base(err)
		}
		return chinasitesrule, nil
	}
	return nil, newError("unknown router rule type: ", rawRule.Type)
}

func parseChinaIPRule(data []byte) (*router.RoutingRule, error) {
	rawRule := new(RouterRule)
	err := json.Unmarshal(data, rawRule)
	if err != nil {
		return nil, newError("invalid router rule").Base(err)
	}
	var chinaIPs geoip.CountryIPRange
	if err := proto.Unmarshal(geoip.ChinaIPs, &chinaIPs); err != nil {
		return nil, newError("invalid china ips").Base(err)
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
		log.Trace(newError("invalid router rule: ", err).AtError())
		return nil, err
	}
	return &router.RoutingRule{
		Tag:    rawRule.OutboundTag,
		Domain: chinaSitesDomains,
	}, nil
}
