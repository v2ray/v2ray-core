package router

import (
	"context"
	"net"

	v2net "v2ray.com/core/common/net"
)

type Rule struct {
	Tag       string
	Condition Condition
}

func (v *Rule) Apply(ctx context.Context) bool {
	return v.Condition.Apply(ctx)
}

func (v *RoutingRule) BuildCondition() (Condition, error) {
	conds := NewConditionChan()

	if len(v.Domain) > 0 {
		anyCond := NewAnyCondition()
		for _, domain := range v.Domain {
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

	if len(v.Cidr) > 0 {
		ipv4Net := v2net.NewIPNet()
		ipv6Cond := NewAnyCondition()
		hasIpv6 := false

		for _, ip := range v.Cidr {
			switch len(ip.Ip) {
			case net.IPv4len:
				ipv4Net.AddIP(ip.Ip, byte(ip.Prefix))
			case net.IPv6len:
				hasIpv6 = true
				matcher, err := NewCIDRMatcher(ip.Ip, ip.Prefix, false)
				if err != nil {
					return nil, err
				}
				ipv6Cond.Add(matcher)
			default:
				return nil, newError("invalid IP length").AtError()
			}
		}

		if !ipv4Net.IsEmpty() && hasIpv6 {
			cond := NewAnyCondition()
			cond.Add(NewIPv4Matcher(ipv4Net, false))
			cond.Add(ipv6Cond)
			conds.Add(cond)
		} else if !ipv4Net.IsEmpty() {
			conds.Add(NewIPv4Matcher(ipv4Net, false))
		} else if hasIpv6 {
			conds.Add(ipv6Cond)
		}
	}

	if v.PortRange != nil {
		conds.Add(NewPortMatcher(*v.PortRange))
	}

	if v.NetworkList != nil {
		conds.Add(NewNetworkMatcher(v.NetworkList))
	}

	if len(v.SourceCidr) > 0 {
		ipv4Net := v2net.NewIPNet()
		ipv6Cond := NewAnyCondition()
		hasIpv6 := false

		for _, ip := range v.SourceCidr {
			switch len(ip.Ip) {
			case net.IPv4len:
				ipv4Net.AddIP(ip.Ip, byte(ip.Prefix))
			case net.IPv6len:
				hasIpv6 = true
				matcher, err := NewCIDRMatcher(ip.Ip, ip.Prefix, true)
				if err != nil {
					return nil, err
				}
				ipv6Cond.Add(matcher)
			default:
				return nil, newError("invalid IP length").AtError()
			}
		}

		if !ipv4Net.IsEmpty() && hasIpv6 {
			cond := NewAnyCondition()
			cond.Add(NewIPv4Matcher(ipv4Net, true))
			cond.Add(ipv6Cond)
			conds.Add(cond)
		} else if !ipv4Net.IsEmpty() {
			conds.Add(NewIPv4Matcher(ipv4Net, true))
		} else if hasIpv6 {
			conds.Add(ipv6Cond)
		}
	}

	if len(v.UserEmail) > 0 {
		conds.Add(NewUserMatcher(v.UserEmail))
	}

	if len(v.InboundTag) > 0 {
		conds.Add(NewInboundTagMatcher(v.InboundTag))
	}

	if conds.Len() == 0 {
		return nil, newError("this rule has no effective fields").AtError()
	}

	return conds, nil
}
