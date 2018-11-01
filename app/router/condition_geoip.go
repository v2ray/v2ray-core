package router

import (
	"sort"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
)

type ipv6 struct {
	a uint64
	b uint64
}

type GeoIPMatcher struct {
	countryCode string
	ip4         []uint32
	prefix4     []uint8
	ip6         []ipv6
	prefix6     []uint8
}

func normalize4(ip uint32, prefix uint8) uint32 {
	return (ip >> (32 - prefix)) << (32 - prefix)
}

func normalize6(ip ipv6, prefix uint8) ipv6 {
	if prefix < 64 {
		ip.a = (ip.a >> (64 - prefix)) << (64 - prefix)
	}

	if prefix <= 64 {
		ip.b = 0
	} else {
		ip.b = (ip.b >> (128 - prefix)) << (128 - prefix)
	}

	return ip
}

func (m *GeoIPMatcher) Init(cidrs []*CIDR) error {
	ip4Count := 0
	ip6Count := 0

	for _, cidr := range cidrs {
		ip := cidr.Ip
		switch len(ip) {
		case 4:
			ip4Count++
		case 16:
			ip6Count++
		default:
			return newError("unexpect ip length: ", len(ip))
		}
	}

	cidrList := CIDRList(cidrs)
	sort.Sort(&cidrList)

	m.ip4 = make([]uint32, 0, ip4Count)
	m.prefix4 = make([]uint8, 0, ip4Count)
	m.ip6 = make([]ipv6, 0, ip6Count)
	m.prefix6 = make([]uint8, 0, ip6Count)

	for _, cidr := range cidrs {
		ip := cidr.Ip
		prefix := uint8(cidr.Prefix)
		switch len(ip) {
		case 4:
			m.ip4 = append(m.ip4, normalize4(serial.BytesToUint32(ip), prefix))
			m.prefix4 = append(m.prefix4, prefix)
		case 16:
			ip6 := ipv6{
				a: serial.BytesToUint64(ip[0:8]),
				b: serial.BytesToUint64(ip[8:16]),
			}
			ip6 = normalize6(ip6, prefix)

			m.ip6 = append(m.ip6, ip6)
			m.prefix6 = append(m.prefix6, prefix)
		}
	}

	return nil
}

func (m *GeoIPMatcher) match4(ip uint32) bool {
	if len(m.ip4) == 0 {
		return false
	}

	if ip < m.ip4[0] {
		return false
	}

	size := uint32(len(m.ip4))
	if ip > m.ip4[size-1] {
		nip := normalize4(ip, m.prefix4[size-1])
		return nip == m.ip4[size-1]
	}

	l := uint32(0)
	r := size - 1
	for l < r-1 {
		x := (l + r) / 2
		if ip < m.ip4[x] {
			r = x
			continue
		}

		nip := normalize4(ip, m.prefix4[x])
		if nip == m.ip4[x] {
			return true
		}

		l = x
	}

	return normalize4(ip, m.prefix4[l]) == m.ip4[l]
}

func less6(a ipv6, b ipv6) bool {
	return a.a < b.a || (a.a == b.a && a.b < b.b)
}

func (m *GeoIPMatcher) match6(ip ipv6) bool {
	if len(m.ip6) == 0 {
		return false
	}

	if less6(ip, m.ip6[0]) {
		return false
	}

	size := uint32(len(m.ip6))
	if less6(m.ip6[size-1], ip) {
		nip := normalize6(ip, m.prefix6[size-1])
		return nip == m.ip6[size-1]
	}

	l := uint32(0)
	r := size - 1
	for l < r-1 {
		x := (l + r) / 2
		if less6(ip, m.ip6[x]) {
			r = x
			continue
		}

		nip := normalize6(ip, m.prefix6[x])
		if nip == m.ip6[x] {
			return true
		}

		l = x
	}

	return normalize6(ip, m.prefix6[l]) == m.ip6[l]
}

// Match returns true if the given ip is included by the GeoIP.
func (m *GeoIPMatcher) Match(ip net.IP) bool {
	switch len(ip) {
	case 4:
		return m.match4(serial.BytesToUint32(ip))
	case 16:
		return m.match6(ipv6{
			a: serial.BytesToUint64(ip[0:8]),
			b: serial.BytesToUint64(ip[8:16]),
		})
	default:
		return false
	}
}

// GeoIPMatcherContainer is a container for GeoIPMatchers. It keeps unique copies of GeoIPMatcher by country code.
type GeoIPMatcherContainer struct {
	matchers []*GeoIPMatcher
}

// Add adds a new GeoIP set into the container.
// If the country code of GeoIP is not empty, GeoIPMatcherContainer will try to find an existing one, instead of adding a new one.
func (c *GeoIPMatcherContainer) Add(geoip *GeoIP) (*GeoIPMatcher, error) {
	if len(geoip.CountryCode) > 0 {
		for _, m := range c.matchers {
			if m.countryCode == geoip.CountryCode {
				return m, nil
			}
		}
	}

	m := &GeoIPMatcher{
		countryCode: geoip.CountryCode,
	}
	if err := m.Init(geoip.Cidr); err != nil {
		return nil, err
	}
	if len(geoip.CountryCode) > 0 {
		c.matchers = append(c.matchers, m)
	}
	return m, nil
}

var (
	globalGeoIPContainer GeoIPMatcherContainer
)
