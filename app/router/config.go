package router

import (
	"errors"
	"net"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type Rule struct {
	Tag       string
	Condition Condition
}

func (this *Rule) Apply(session *proxy.SessionInfo) bool {
	return this.Condition.Apply(session)
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

	if len(this.Cidr) > 0 {
		ipv4Net := v2net.NewIPNet()
		ipv6Cond := NewAnyCondition()
		hasIpv6 := false

		for _, ip := range this.Cidr {
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
				return nil, errors.New("Router: Invalid IP length.")
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

	if this.PortRange != nil {
		conds.Add(NewPortMatcher(*this.PortRange))
	}

	if this.NetworkList != nil {
		conds.Add(NewNetworkMatcher(this.NetworkList))
	}

	if len(this.SourceCidr) > 0 {
		ipv4Net := v2net.NewIPNet()
		ipv6Cond := NewAnyCondition()
		hasIpv6 := false

		for _, ip := range this.SourceCidr {
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
				return nil, errors.New("Router: Invalid IP length.")
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

	if len(this.UserEmail) > 0 {
		conds.Add(NewUserMatcher(this.UserEmail))
	}

	if conds.Len() == 0 {
		return nil, errors.New("Router: This rule has no effective fields.")
	}

	return conds, nil
}
