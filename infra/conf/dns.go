package conf

import (
	"encoding/json"
	"sort"
	"strings"

	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common/net"
)

type NameServerConfig struct {
	Address *Address
	Port    uint16
	Domains []string
}

func (c *NameServerConfig) UnmarshalJSON(data []byte) error {
	var address Address
	if err := json.Unmarshal(data, &address); err == nil {
		c.Address = &address
		c.Port = 53
		return nil
	}

	var advanced struct {
		Address *Address `json:"address"`
		Port    uint16   `json:"port"`
		Domains []string `json:"domains"`
	}
	if err := json.Unmarshal(data, &advanced); err == nil {
		c.Address = advanced.Address
		c.Port = advanced.Port
		c.Domains = advanced.Domains
		return nil
	}

	return newError("failed to parse name server: ", string(data))
}

func toDomainMatchingType(t router.Domain_Type) dns.DomainMatchingType {
	switch t {
	case router.Domain_Domain:
		return dns.DomainMatchingType_Subdomain
	case router.Domain_Full:
		return dns.DomainMatchingType_Full
	case router.Domain_Plain:
		return dns.DomainMatchingType_Keyword
	case router.Domain_Regex:
		return dns.DomainMatchingType_Regex
	default:
		panic("unknown domain type")
	}
}

func (c *NameServerConfig) Build() (*dns.NameServer, error) {
	if c.Address == nil {
		return nil, newError("NameServer address is not specified.")
	}

	var domains []*dns.NameServer_PriorityDomain

	for _, d := range c.Domains {
		parsedDomain, err := parseDomainRule(d)
		if err != nil {
			return nil, newError("invalid domain rule: ", d).Base(err)
		}

		for _, pd := range parsedDomain {
			domains = append(domains, &dns.NameServer_PriorityDomain{
				Type:   toDomainMatchingType(pd.Type),
				Domain: pd.Value,
			})
		}
	}

	return &dns.NameServer{
		Address: &net.Endpoint{
			Network: net.Network_UDP,
			Address: c.Address.Build(),
			Port:    uint32(c.Port),
		},
		PrioritizedDomain: domains,
	}, nil
}

var typeMap = map[router.Domain_Type]dns.DomainMatchingType{
	router.Domain_Full:   dns.DomainMatchingType_Full,
	router.Domain_Domain: dns.DomainMatchingType_Subdomain,
	router.Domain_Plain:  dns.DomainMatchingType_Keyword,
	router.Domain_Regex:  dns.DomainMatchingType_Regex,
}

// DnsConfig is a JSON serializable object for dns.Config.
type DnsConfig struct {
	Servers  []*NameServerConfig `json:"servers"`
	Hosts    map[string]*Address `json:"hosts"`
	ClientIP *Address            `json:"clientIp"`
	Tag      string              `json:"tag"`
}

func getHostMapping(addr *Address) *dns.Config_HostMapping {
	if addr.Family().IsIP() {
		return &dns.Config_HostMapping{
			Ip: [][]byte{[]byte(addr.IP())},
		}
	} else {
		return &dns.Config_HostMapping{
			ProxiedDomain: addr.Domain(),
		}
	}
}

// Build implements Buildable
func (c *DnsConfig) Build() (*dns.Config, error) {
	config := &dns.Config{
		Tag: c.Tag,
	}

	if c.ClientIP != nil {
		if !c.ClientIP.Family().IsIP() {
			return nil, newError("not an IP address:", c.ClientIP.String())
		}
		config.ClientIp = []byte(c.ClientIP.IP())
	}

	for _, server := range c.Servers {
		ns, err := server.Build()
		if err != nil {
			return nil, newError("failed to build name server").Base(err)
		}
		config.NameServer = append(config.NameServer, ns)
	}

	if c.Hosts != nil && len(c.Hosts) > 0 {
		domains := make([]string, 0, len(c.Hosts))
		for domain := range c.Hosts {
			domains = append(domains, domain)
		}
		sort.Strings(domains)
		for _, domain := range domains {
			addr := c.Hosts[domain]
			var mappings []*dns.Config_HostMapping
			if strings.HasPrefix(domain, "domain:") {
				mapping := getHostMapping(addr)
				mapping.Type = dns.DomainMatchingType_Subdomain
				mapping.Domain = domain[7:]

				mappings = append(mappings, mapping)
			} else if strings.HasPrefix(domain, "geosite:") {
				domains, err := loadGeositeWithAttr("geosite.dat", strings.ToUpper(domain[8:]))
				if err != nil {
					return nil, newError("invalid geosite settings: ", domain).Base(err)
				}
				for _, d := range domains {
					mapping := getHostMapping(addr)
					mapping.Type = typeMap[d.Type]
					mapping.Domain = d.Value

					mappings = append(mappings, mapping)
				}
			} else if strings.HasPrefix(domain, "regexp:") {
				mapping := getHostMapping(addr)
				mapping.Type = dns.DomainMatchingType_Regex
				mapping.Domain = domain[7:]

				mappings = append(mappings, mapping)
			} else if strings.HasPrefix(domain, "keyword:") {
				mapping := getHostMapping(addr)
				mapping.Type = dns.DomainMatchingType_Keyword
				mapping.Domain = domain[8:]

				mappings = append(mappings, mapping)
			} else if strings.HasPrefix(domain, "full:") {
				mapping := getHostMapping(addr)
				mapping.Type = dns.DomainMatchingType_Full
				mapping.Domain = domain[5:]

				mappings = append(mappings, mapping)
			} else {
				mapping := getHostMapping(addr)
				mapping.Type = dns.DomainMatchingType_Full
				mapping.Domain = domain

				mappings = append(mappings, mapping)
			}

			config.StaticHosts = append(config.StaticHosts, mappings...)
		}
	}

	return config, nil
}
