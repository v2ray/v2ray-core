package router

import (
	"net"
	"regexp"
	"strings"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type Condition interface {
	Apply(session *proxy.SessionInfo) bool
}

type ConditionChan []Condition

func NewConditionChan() *ConditionChan {
	var condChan ConditionChan = make([]Condition, 0, 8)
	return &condChan
}

func (v *ConditionChan) Add(cond Condition) *ConditionChan {
	*v = append(*v, cond)
	return v
}

func (v *ConditionChan) Apply(session *proxy.SessionInfo) bool {
	for _, cond := range *v {
		if !cond.Apply(session) {
			return false
		}
	}
	return true
}

func (v *ConditionChan) Len() int {
	return len(*v)
}

type AnyCondition []Condition

func NewAnyCondition() *AnyCondition {
	var anyCond AnyCondition = make([]Condition, 0, 8)
	return &anyCond
}

func (v *AnyCondition) Add(cond Condition) *AnyCondition {
	*v = append(*v, cond)
	return v
}

func (v *AnyCondition) Apply(session *proxy.SessionInfo) bool {
	for _, cond := range *v {
		if cond.Apply(session) {
			return true
		}
	}
	return false
}

func (v *AnyCondition) Len() int {
	return len(*v)
}

type PlainDomainMatcher struct {
	pattern string
}

func NewPlainDomainMatcher(pattern string) *PlainDomainMatcher {
	return &PlainDomainMatcher{
		pattern: pattern,
	}
}

func (v *PlainDomainMatcher) Apply(session *proxy.SessionInfo) bool {
	dest := session.Destination
	if !dest.Address.Family().IsDomain() {
		return false
	}
	domain := dest.Address.Domain()
	return strings.Contains(domain, v.pattern)
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

func (v *RegexpDomainMatcher) Apply(session *proxy.SessionInfo) bool {
	dest := session.Destination
	if !dest.Address.Family().IsDomain() {
		return false
	}
	domain := dest.Address.Domain()
	return v.pattern.MatchString(strings.ToLower(domain))
}

type CIDRMatcher struct {
	cidr     *net.IPNet
	onSource bool
}

func NewCIDRMatcher(ip []byte, mask uint32, onSource bool) (*CIDRMatcher, error) {
	cidr := &net.IPNet{
		IP:   net.IP(ip),
		Mask: net.CIDRMask(int(mask), len(ip)),
	}
	return &CIDRMatcher{
		cidr:     cidr,
		onSource: onSource,
	}, nil
}

func (v *CIDRMatcher) Apply(session *proxy.SessionInfo) bool {
	dest := session.Destination
	if v.onSource {
		dest = session.Source
	}
	if !dest.Address.Family().Either(v2net.AddressFamilyIPv4, v2net.AddressFamilyIPv6) {
		return false
	}
	return v.cidr.Contains(dest.Address.IP())
}

type IPv4Matcher struct {
	ipv4net  *v2net.IPNet
	onSource bool
}

func NewIPv4Matcher(ipnet *v2net.IPNet, onSource bool) *IPv4Matcher {
	return &IPv4Matcher{
		ipv4net:  ipnet,
		onSource: onSource,
	}
}

func (v *IPv4Matcher) Apply(session *proxy.SessionInfo) bool {
	dest := session.Destination
	if v.onSource {
		dest = session.Source
	}
	if !dest.Address.Family().Either(v2net.AddressFamilyIPv4) {
		return false
	}
	return v.ipv4net.Contains(dest.Address.IP())
}

type PortMatcher struct {
	port v2net.PortRange
}

func NewPortMatcher(portRange v2net.PortRange) *PortMatcher {
	return &PortMatcher{
		port: portRange,
	}
}

func (v *PortMatcher) Apply(session *proxy.SessionInfo) bool {
	return v.port.Contains(session.Destination.Port)
}

type NetworkMatcher struct {
	network *v2net.NetworkList
}

func NewNetworkMatcher(network *v2net.NetworkList) *NetworkMatcher {
	return &NetworkMatcher{
		network: network,
	}
}

func (v *NetworkMatcher) Apply(session *proxy.SessionInfo) bool {
	return v.network.HasNetwork(session.Destination.Network)
}

type UserMatcher struct {
	user []string
}

func NewUserMatcher(users []string) *UserMatcher {
	return &UserMatcher{
		user: users,
	}
}

func (v *UserMatcher) Apply(session *proxy.SessionInfo) bool {
	if session.User == nil {
		return false
	}
	for _, u := range v.user {
		if u == session.User.Email {
			return true
		}
	}
	return false
}

type InboundTagMatcher struct {
	tags []string
}

func NewInboundTagMatcher(tags []string) *InboundTagMatcher {
	return &InboundTagMatcher{
		tags: tags,
	}
}

func (v *InboundTagMatcher) Apply(session *proxy.SessionInfo) bool {
	if session.Inbound == nil || len(session.Inbound.Tag) == 0 {
		return false
	}

	for _, t := range v.tags {
		if t == session.Inbound.Tag {
			return true
		}
	}
	return false
}
