package router

import (
	"errors"
	"net"

	v2net "v2ray.com/core/common/net"
)

type Rule struct {
	Tag       string
	Condition Condition
}

func (this *Rule) Apply(dest v2net.Destination) bool {
	return this.Condition.Apply(dest)
}

func (this *RoutingRule) BuildCondition() (Condition, error) {
	conds := NewConditionChan()

	if len(this.Domain) > 0 {
		anyCond := NewAnyCondition()
		for _, domain := range this.Domain {
			if domain.Type == Domain_Plain {
				anyCond.Add(NewPlainDomainMatcher(domain.Value))
			} else {
				matcher, err := NewRegexpDomainMatcher(domain.Value)
				if err != nil {
					return nil, err
				}
				anyCond.Add(matcher)
			}
		}
		conds.Add(anyCond)
	}

	if len(this.Ip) > 0 {
		ipv4Net := make(map[uint32]byte)
		ipv6Cond := NewAnyCondition()
		hasIpv6 := false

		for _, ip := range this.Ip {
			switch len(ip.Ip) {
			case net.IPv4len:
				k := (uint32(ip.Ip[0]) << 24) + (uint32(ip.Ip[1]) << 16) + (uint32(ip.Ip[2]) << 8) + uint32(ip.Ip[3])
				ipv4Net[k] = byte(32 - ip.UnmatchingBits)
			case net.IPv6len:
				hasIpv6 = true
				matcher, err := NewCIDRMatcher(ip.Ip, uint32(32)-ip.UnmatchingBits)
				if err != nil {
					return nil, err
				}
				ipv6Cond.Add(matcher)
			default:
				return nil, errors.New("Router: Invalid IP length.")
			}
		}

		if len(ipv4Net) > 0 && hasIpv6 {
			cond := NewAnyCondition()
			cond.Add(NewIPv4Matcher(v2net.NewIPNetInitialValue(ipv4Net)))
			cond.Add(ipv6Cond)
			conds.Add(cond)
		} else if len(ipv4Net) > 0 {
			conds.Add(NewIPv4Matcher(v2net.NewIPNetInitialValue(ipv4Net)))
		} else if hasIpv6 {
			conds.Add(ipv6Cond)
		}
	}

	if this.PortRange != nil {
		conds.Add(NewPortMatcher(*this.PortRange))
	}

	if this.NetworkList != nil {
		conds.Add(NewNetworkMatcher(this.NetworkList))
	}

	if conds.Len() == 0 {
		return nil, errors.New("Router: This rule has no effective fields.")
	}

	return conds, nil
}
