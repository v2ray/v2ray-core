// +build !confonly

package router

import (
	"strings"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
)

type Condition interface {
	Apply(ctx *Context) bool
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

func (v *ConditionChan) Apply(ctx *Context) bool {
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

func (m *DomainMatcher) Apply(ctx *Context) bool {
	if ctx.Outbound == nil || !ctx.Outbound.Target.IsValid() {
		return false
	}
	dest := ctx.Outbound.Target
	if !dest.Address.Family().IsDomain() {
		return false
	}
	return m.ApplyDomain(dest.Address.Domain())
}

func getIPsFromSource(ctx *Context) []net.IP {
	if ctx.Inbound == nil || !ctx.Inbound.Source.IsValid() {
		return nil
	}
	dest := ctx.Inbound.Source
	if dest.Address.Family().IsDomain() {
		return nil
	}

	return []net.IP{dest.Address.IP()}
}

func getIPsFromTarget(ctx *Context) []net.IP {
	return ctx.GetTargetIPs()
}

type MultiGeoIPMatcher struct {
	matchers []*GeoIPMatcher
	ipFunc   func(*Context) []net.IP
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

	matcher := &MultiGeoIPMatcher{
		matchers: matchers,
	}

	if onSource {
		matcher.ipFunc = getIPsFromSource
	} else {
		matcher.ipFunc = getIPsFromTarget
	}

	return matcher, nil
}

func (m *MultiGeoIPMatcher) Apply(ctx *Context) bool {
	ips := m.ipFunc(ctx)

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
	port net.MemoryPortList
}

func NewPortMatcher(list *net.PortList) *PortMatcher {
	return &PortMatcher{
		port: net.PortListFromProto(list),
	}
}

func (v *PortMatcher) Apply(ctx *Context) bool {
	if ctx.Outbound == nil || !ctx.Outbound.Target.IsValid() {
		return false
	}
	return v.port.Contains(ctx.Outbound.Target.Port)
}

type NetworkMatcher struct {
	list [8]bool
}

func NewNetworkMatcher(network []net.Network) NetworkMatcher {
	var matcher NetworkMatcher
	for _, n := range network {
		matcher.list[int(n)] = true
	}
	return matcher
}

func (v NetworkMatcher) Apply(ctx *Context) bool {
	if ctx.Outbound == nil || !ctx.Outbound.Target.IsValid() {
		return false
	}
	return v.list[int(ctx.Outbound.Target.Network)]
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

func (v *UserMatcher) Apply(ctx *Context) bool {
	if ctx.Inbound == nil {
		return false
	}

	user := ctx.Inbound.User
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

func (v *InboundTagMatcher) Apply(ctx *Context) bool {
	if ctx.Inbound == nil || len(ctx.Inbound.Tag) == 0 {
		return false
	}
	tag := ctx.Inbound.Tag
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

func (m *ProtocolMatcher) Apply(ctx *Context) bool {
	if ctx.Content == nil {
		return false
	}

	protocol := ctx.Content.Protocol
	for _, p := range m.protocols {
		if strings.HasPrefix(protocol, p) {
			return true
		}
	}

	return false
}
