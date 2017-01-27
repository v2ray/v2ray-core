package router

import (
	"context"
	"net"
	"regexp"
	"strings"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy"
)

type Condition interface {
	Apply(ctx context.Context) bool
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

func (v *ConditionChan) Apply(ctx context.Context) bool {
	for _, cond := range *v {
		if !cond.Apply(ctx) {
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

func (v *AnyCondition) Apply(ctx context.Context) bool {
	for _, cond := range *v {
		if cond.Apply(ctx) {
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

func (v *PlainDomainMatcher) Apply(ctx context.Context) bool {
	dest := proxy.DestinationFromContext(ctx)
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

func (v *RegexpDomainMatcher) Apply(ctx context.Context) bool {
	dest := proxy.DestinationFromContext(ctx)
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

func (v *CIDRMatcher) Apply(ctx context.Context) bool {
	ips := make([]net.IP, 4)
	if resolveIPs, ok := proxy.ResolvedIPsFromContext(ctx); ok {
		for _, rip := range resolveIPs {
			ips = append(ips, rip.IP())
		}
	}

	var dest v2net.Destination
	if v.onSource {
		dest = proxy.SourceFromContext(ctx)
	} else {
		dest = proxy.DestinationFromContext(ctx)
	}

	if dest.IsValid() && dest.Address.Family().Either(v2net.AddressFamilyIPv4, v2net.AddressFamilyIPv6) {
		ips = append(ips, dest.Address.IP())
	}

	for _, ip := range ips {
		if v.cidr.Contains(ip) {
			return true
		}
	}
	return false
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

func (v *IPv4Matcher) Apply(ctx context.Context) bool {
	ips := make([]net.IP, 4)
	if resolveIPs, ok := proxy.ResolvedIPsFromContext(ctx); ok {
		for _, rip := range resolveIPs {
			ips = append(ips, rip.IP())
		}
	}

	var dest v2net.Destination
	if v.onSource {
		dest = proxy.SourceFromContext(ctx)
	} else {
		dest = proxy.DestinationFromContext(ctx)
	}

	if dest.IsValid() && dest.Address.Family().Either(v2net.AddressFamilyIPv4) {
		ips = append(ips, dest.Address.IP())
	}

	for _, ip := range ips {
		if v.ipv4net.Contains(ip) {
			return true
		}
	}
	return false
}

type PortMatcher struct {
	port v2net.PortRange
}

func NewPortMatcher(portRange v2net.PortRange) *PortMatcher {
	return &PortMatcher{
		port: portRange,
	}
}

func (v *PortMatcher) Apply(ctx context.Context) bool {
	dest := proxy.DestinationFromContext(ctx)
	return v.port.Contains(dest.Port)
}

type NetworkMatcher struct {
	network *v2net.NetworkList
}

func NewNetworkMatcher(network *v2net.NetworkList) *NetworkMatcher {
	return &NetworkMatcher{
		network: network,
	}
}

func (v *NetworkMatcher) Apply(ctx context.Context) bool {
	dest := proxy.DestinationFromContext(ctx)
	return v.network.HasNetwork(dest.Network)
}

type UserMatcher struct {
	user []string
}

func NewUserMatcher(users []string) *UserMatcher {
	return &UserMatcher{
		user: users,
	}
}

func (v *UserMatcher) Apply(ctx context.Context) bool {
	user := protocol.UserFromContext(ctx)
	if user == nil {
		return false
	}
	for _, u := range v.user {
		if u == user.Email {
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

func (v *InboundTagMatcher) Apply(ctx context.Context) bool {
	tag := proxy.InboundTagFromContext(ctx)
	if len(tag) == 0 {
		return false
	}

	for _, t := range v.tags {
		if t == tag {
			return true
		}
	}
	return false
}
