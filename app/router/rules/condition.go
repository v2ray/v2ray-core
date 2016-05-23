package rules

import (
	"net"
	"regexp"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
)

type Condition interface {
	Apply(dest v2net.Destination) bool
}

type ConditionChan []Condition

func NewConditionChan() *ConditionChan {
	var condChan ConditionChan = make([]Condition, 0, 8)
	return &condChan
}

func (this *ConditionChan) Add(cond Condition) *ConditionChan {
	*this = append(*this, cond)
	return this
}

func (this *ConditionChan) Apply(dest v2net.Destination) bool {
	for _, cond := range *this {
		if !cond.Apply(dest) {
			return false
		}
	}
	return true
}

func (this *ConditionChan) Len() int {
	return len(*this)
}

type AnyCondition []Condition

func NewAnyCondition() *AnyCondition {
	var anyCond AnyCondition = make([]Condition, 0, 8)
	return &anyCond
}

func (this *AnyCondition) Add(cond Condition) *AnyCondition {
	*this = append(*this, cond)
	return this
}

func (this *AnyCondition) Apply(dest v2net.Destination) bool {
	for _, cond := range *this {
		if cond.Apply(dest) {
			return true
		}
	}
	return false
}

func (this *AnyCondition) Len() int {
	return len(*this)
}

type PlainDomainMatcher struct {
	pattern serial.StringT
}

func NewPlainDomainMatcher(pattern string) *PlainDomainMatcher {
	return &PlainDomainMatcher{
		pattern: serial.StringT(pattern),
	}
}

func (this *PlainDomainMatcher) Apply(dest v2net.Destination) bool {
	if !dest.Address().IsDomain() {
		return false
	}
	domain := serial.StringT(dest.Address().Domain())
	return domain.Contains(this.pattern)
}

type RegexpDomainMatcher struct {
	pattern *regexp.Regexp
}

func NewRegexpDomainMatcher(pattern string) (*RegexpDomainMatcher, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &RegexpDomainMatcher{
		pattern: r,
	}, nil
}

func (this *RegexpDomainMatcher) Apply(dest v2net.Destination) bool {
	if !dest.Address().IsDomain() {
		return false
	}
	domain := serial.StringT(dest.Address().Domain())
	return this.pattern.MatchString(domain.ToLower().String())
}

type CIDRMatcher struct {
	cidr *net.IPNet
}

func NewCIDRMatcher(ipnet string) (*CIDRMatcher, error) {
	_, cidr, err := net.ParseCIDR(ipnet)
	if err != nil {
		return nil, err
	}
	return &CIDRMatcher{
		cidr: cidr,
	}, nil
}

func (this *CIDRMatcher) Apply(dest v2net.Destination) bool {
	if !dest.Address().IsIPv4() && !dest.Address().IsIPv6() {
		return false
	}
	return this.cidr.Contains(dest.Address().IP())
}

type IPv4Matcher struct {
	ipv4net *v2net.IPNet
}

func NewIPv4Matcher(ipnet *v2net.IPNet) *IPv4Matcher {
	return &IPv4Matcher{
		ipv4net: ipnet,
	}
}

func (this *IPv4Matcher) Apply(dest v2net.Destination) bool {
	if !dest.Address().IsIPv4() {
		return false
	}
	return this.ipv4net.Contains(dest.Address().IP())
}

type PortMatcher struct {
	port v2net.PortRange
}

func NewPortMatcher(portRange v2net.PortRange) *PortMatcher {
	return &PortMatcher{
		port: portRange,
	}
}

func (this *PortMatcher) Apply(dest v2net.Destination) bool {
	return this.port.Contains(dest.Port())
}

type NetworkMatcher struct {
	network *v2net.NetworkList
}

func NewNetworkMatcher(network *v2net.NetworkList) *NetworkMatcher {
	return &NetworkMatcher{
		network: network,
	}
}

func (this *NetworkMatcher) Apply(dest v2net.Destination) bool {
	return this.network.HasNetwork(dest.Network())
}
