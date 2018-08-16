package dns

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
)

type StaticHosts struct {
	ips      map[uint32][]net.IP
	matchers *strmatcher.MatcherGroup
}

var typeMap = map[Config_HostMapping_Type]strmatcher.Type{
	Config_HostMapping_Full:      strmatcher.Full,
	Config_HostMapping_SubDomain: strmatcher.Domain,
}

func NewStaticHosts(hosts []*Config_HostMapping, legacy map[string]*net.IPOrDomain) (*StaticHosts, error) {
	g := strmatcher.NewMatcherGroup()
	sh := &StaticHosts{
		ips:      make(map[uint32][]net.IP),
		matchers: g,
	}

	if legacy != nil {
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
		strMType, f := typeMap[mapping.Type]
		if !f {
			return nil, newError("unknown mapping type", mapping.Type).AtWarning()
		}
		matcher, err := strMType.New(mapping.Domain)
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
