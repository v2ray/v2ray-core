package dns

import (
	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
)

type StaticHosts struct {
	ips      map[uint32][]net.IP
	matchers *strmatcher.MatcherGroup
}

var typeMap = map[DomainMatchingType]strmatcher.Type{
	DomainMatchingType_Full:      strmatcher.Full,
	DomainMatchingType_Subdomain: strmatcher.Domain,
	DomainMatchingType_Keyword:   strmatcher.Substr,
	DomainMatchingType_Regex:     strmatcher.Regex,
}

func toStrMatcher(t DomainMatchingType, domain string) (strmatcher.Matcher, error) {
	strMType, f := typeMap[t]
	if !f {
		return nil, newError("unknown mapping type", t).AtWarning()
	}
	matcher, err := strMType.New(domain)
	if err != nil {
		return nil, newError("failed to create str matcher").Base(err)
	}
	return matcher, nil
}

func NewStaticHosts(hosts []*Config_HostMapping, legacy map[string]*net.IPOrDomain) (*StaticHosts, error) {
	g := new(strmatcher.MatcherGroup)
	sh := &StaticHosts{
		ips:      make(map[uint32][]net.IP),
		matchers: g,
	}

	if legacy != nil {
		core.PrintDeprecatedFeatureWarning("simple host mapping")

		for domain, ip := range legacy {
			matcher, err := strmatcher.Full.New(domain)
			common.Must(err)
			id := g.Add(matcher)

			address := ip.AsAddress()
			if address.Family().IsDomain() {
				return nil, newError("ignoring domain address in static hosts: ", address.Domain()).AtWarning()
			}

			sh.ips[id] = []net.IP{address.IP()}
		}
	}

	for _, mapping := range hosts {
		matcher, err := toStrMatcher(mapping.Type, mapping.Domain)
		if err != nil {
			return nil, newError("failed to create domain matcher").Base(err)
		}
		id := g.Add(matcher)
		ips := make([]net.IP, len(mapping.Ip))
		for idx, ip := range mapping.Ip {
			ips[idx] = net.IP(ip)
		}
		sh.ips[id] = ips
	}

	return sh, nil
}

func (h *StaticHosts) LookupIP(domain string) []net.IP {
	id := h.matchers.Match(domain)
	if id == 0 {
		return nil
	}
	return h.ips[id]
}
