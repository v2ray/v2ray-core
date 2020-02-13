// +build !confonly

package dns

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
	"v2ray.com/core/features"
)

// StaticHosts represents static domain-ip mapping in DNS server.
type StaticHosts struct {
	ips      [][]net.Address
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

var typeMapper = map[DomainMatchingType]string{
	DomainMatchingType_Full:      "f",
	DomainMatchingType_Subdomain: "d",
	DomainMatchingType_Keyword:   "k",
	DomainMatchingType_Regex:     "r",
}

// NewStaticHosts creates a new StaticHosts instance.
func NewStaticHosts(hosts []*Config_HostMapping, legacy map[string]*net.IPOrDomain, externalRules map[string][]string) (*StaticHosts, error) {
	g := new(strmatcher.MatcherGroup)
	sh := &StaticHosts{
		ips:      make([][]net.Address, len(hosts)+len(legacy)+16),
		matchers: g,
	}

	if legacy != nil {
		features.PrintDeprecatedFeatureWarning("simple host mapping")

		for domain, ip := range legacy {
			id, err := g.ParsePattern("f"+domain, make(map[string][]string))
			common.Must(err)

			address := ip.AsAddress()
			if address.Family().IsDomain() {
				return nil, newError("invalid domain address in static hosts: ", address.Domain()).AtWarning()
			}

			sh.ips[id] = []net.Address{address}
		}
	}

	for _, mapping := range hosts {
		if mapping.Type != DomainMatchingType_New {
			mapping.Domain = typeMapper[mapping.Type] + mapping.Domain
		}
		id, err := g.ParsePattern(mapping.Domain, externalRules)
		if err != nil {
			return nil, newError("failed to create domain matcher").Base(err)
		}
		ips := make([]net.Address, 0, len(mapping.Ip)+1)
		if len(mapping.Ip) > 0 {
			for _, ip := range mapping.Ip {
				addr := net.IPAddress(ip)
				if addr == nil {
					return nil, newError("invalid IP address in static hosts: ", ip).AtWarning()
				}
				ips = append(ips, addr)
			}
		} else if len(mapping.ProxiedDomain) > 0 {
			ips = append(ips, net.DomainAddress(mapping.ProxiedDomain))
		} else {
			return nil, newError("neither IP address nor proxied domain specified for domain: ", mapping.Domain).AtWarning()
		}

		// Special handling for localhost IPv6. This is a dirty workaround as JSON config supports only single IP mapping.
		if len(ips) == 1 && ips[0] == net.LocalHostIP {
			ips = append(ips, net.LocalHostIPv6)
		}

		sh.ips[id] = ips
	}

	return sh, nil
}

func filterIP(ips []net.Address, option IPOption) []net.Address {
	filtered := make([]net.Address, 0, len(ips))
	for _, ip := range ips {
		if (ip.Family().IsIPv4() && option.IPv4Enable) || (ip.Family().IsIPv6() && option.IPv6Enable) {
			filtered = append(filtered, ip)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return filtered
}

// LookupIP returns IP address for the given domain, if exists in this StaticHosts.
func (h *StaticHosts) LookupIP(domain string, option IPOption) []net.Address {
	id := h.matchers.Match(domain)
	if id == 0 {
		return nil
	}
	ips := h.ips[id]
	if len(ips) == 1 && ips[0].Family().IsDomain() {
		return ips
	}
	return filterIP(ips, option)
}
