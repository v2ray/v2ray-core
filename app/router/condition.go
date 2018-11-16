package router

import (
	"context"
	"strings"

	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/strmatcher"
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

var matcherTypeMap = map[Domain_Type]strmatcher.Type{
	Domain_Plain:  strmatcher.Substr,
	Domain_Regex:  strmatcher.Regex,
	Domain_Domain: strmatcher.Domain,
	Domain_Full:   strmatcher.Full,
}

func domainToMatcher(domain *Domain) (strmatcher.Matcher, error) {
	matcherType, f := matcherTypeMap[domain.Type]
	if !f {
		return nil, newError("unsupported domain type", domain.Type)
	}

	matcher, err := matcherType.New(domain.Value)
	if err != nil {
		return nil, newError("failed to create domain matcher").Base(err)
	}

	return matcher, nil
}

type DomainMatcher struct {
	matchers strmatcher.IndexMatcher
}

func NewDomainMatcher(domains []*Domain) (*DomainMatcher, error) {
	g := new(strmatcher.MatcherGroup)
	for _, d := range domains {
		m, err := domainToMatcher(d)
		if err != nil {
			return nil, err
		}
		g.Add(m)
	}

	return &DomainMatcher{
		matchers: g,
	}, nil
}

func (m *DomainMatcher) ApplyDomain(domain string) bool {
	return m.matchers.Match(domain) > 0
}

func (m *DomainMatcher) Apply(ctx context.Context) bool {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return false
	}
	dest := outbound.Target
	if !dest.Address.Family().IsDomain() {
		return false
	}
	return m.ApplyDomain(dest.Address.Domain())
}

func sourceFromContext(ctx context.Context) net.Destination {
	inbound := session.InboundFromContext(ctx)
	if inbound == nil {
		return net.Destination{}
	}
	return inbound.Source
}

func targetFromContent(ctx context.Context) net.Destination {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil {
		return net.Destination{}
	}
	return outbound.Target
}

type MultiGeoIPMatcher struct {
	matchers []*GeoIPMatcher
	destFunc func(context.Context) net.Destination
}

func NewMultiGeoIPMatcher(geoips []*GeoIP, onSource bool) (*MultiGeoIPMatcher, error) {
	var matchers []*GeoIPMatcher
	for _, geoip := range geoips {
		matcher, err := globalGeoIPContainer.Add(geoip)
		if err != nil {
			return nil, err
		}
		matchers = append(matchers, matcher)
	}

	var destFunc func(context.Context) net.Destination
	if onSource {
		destFunc = sourceFromContext
	} else {
		destFunc = targetFromContent
	}

	return &MultiGeoIPMatcher{
		matchers: matchers,
		destFunc: destFunc,
	}, nil
}

func (m *MultiGeoIPMatcher) Apply(ctx context.Context) bool {
	ips := make([]net.IP, 0, 4)

	dest := m.destFunc(ctx)

	if dest.IsValid() && dest.Address.Family().IsIP() {
		ips = append(ips, dest.Address.IP())
	} else if resolver, ok := ResolvedIPsFromContext(ctx); ok {
		resolvedIPs := resolver.Resolve()
		for _, rip := range resolvedIPs {
			ips = append(ips, rip.IP())
		}
	}

	for _, ip := range ips {
		for _, matcher := range m.matchers {
			if matcher.Match(ip) {
				return true
			}
		}
	}
	return false
}

type PortMatcher struct {
	port net.PortRange
}

func NewPortMatcher(portRange net.PortRange) *PortMatcher {
	return &PortMatcher{
		port: portRange,
	}
}

func (v *PortMatcher) Apply(ctx context.Context) bool {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return false
	}
	return v.port.Contains(outbound.Target.Port)
}

type NetworkMatcher struct {
	network *net.NetworkList
}

func NewNetworkMatcher(network *net.NetworkList) *NetworkMatcher {
	return &NetworkMatcher{
		network: network,
	}
}

func (v *NetworkMatcher) Apply(ctx context.Context) bool {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return false
	}
	return v.network.HasNetwork(outbound.Target.Network)
}

type UserMatcher struct {
	user []string
}

func NewUserMatcher(users []string) *UserMatcher {
	usersCopy := make([]string, 0, len(users))
	for _, user := range users {
		if len(user) > 0 {
			usersCopy = append(usersCopy, user)
		}
	}
	return &UserMatcher{
		user: usersCopy,
	}
}

func (v *UserMatcher) Apply(ctx context.Context) bool {
	inbound := session.InboundFromContext(ctx)
	if inbound == nil {
		return false
	}

	user := inbound.User
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
	tagsCopy := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag) > 0 {
			tagsCopy = append(tagsCopy, tag)
		}
	}
	return &InboundTagMatcher{
		tags: tagsCopy,
	}
}

func (v *InboundTagMatcher) Apply(ctx context.Context) bool {
	inbound := session.InboundFromContext(ctx)
	if inbound == nil || len(inbound.Tag) == 0 {
		return false
	}
	tag := inbound.Tag
	for _, t := range v.tags {
		if t == tag {
			return true
		}
	}
	return false
}

type ProtocolMatcher struct {
	protocols []string
}

func NewProtocolMatcher(protocols []string) *ProtocolMatcher {
	pCopy := make([]string, 0, len(protocols))

	for _, p := range protocols {
		if len(p) > 0 {
			pCopy = append(pCopy, p)
		}
	}

	return &ProtocolMatcher{
		protocols: pCopy,
	}
}

func (m *ProtocolMatcher) Apply(ctx context.Context) bool {
	result := dispatcher.SniffingResultFromContext(ctx)

	if result == nil {
		return false
	}

	protocol := result.Protocol()
	for _, p := range m.protocols {
		if strings.HasPrefix(protocol, p) {
			return true
		}
	}

	return false
}
