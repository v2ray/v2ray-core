package conf

import (
	"encoding/json"

	"v2ray.com/core/app/dns"
	"v2ray.com/core/common/net"
)

type NameServerConfig struct {
	Address   *Address
	Port      uint16
	Domains   []string
	ExpectIPs StringList
}

func (c *NameServerConfig) UnmarshalJSON(data []byte) error {
	var address Address
	if err := json.Unmarshal(data, &address); err == nil {
		c.Address = &address
		return nil
	}

	var advanced struct {
		Address   *Address   `json:"address"`
		Port      uint16     `json:"port"`
		Domains   []string   `json:"domains"`
		ExpectIPs StringList `json:"expectIps"`
	}
	if err := json.Unmarshal(data, &advanced); err == nil {
		c.Address = advanced.Address
		c.Port = advanced.Port
		c.Domains = advanced.Domains
		c.ExpectIPs = advanced.ExpectIPs
		return nil
	}

	return newError("failed to parse name server: ", string(data))
}

// FakeIPConfig contains configurations for fake IP function
type FakeIPConfig struct {
	FakeRules []string
	FakeNet   string
}

// UnmarshalJSON is an implemention for unmarshal json data
func (c *FakeIPConfig) UnmarshalJSON(data []byte) error {
	var advanced struct {
		FakeRules []string `json:"fakeRules"`
		FakeNet   string   `json:"fakeNet"`
	}

	if err := json.Unmarshal(data, &advanced); err == nil {
		c.FakeRules = advanced.FakeRules
		c.FakeNet = advanced.FakeNet
		return nil
	}

	return newError("failed to parse fake config: ", string(data))
}

var externalRules = make(map[string][]string)

func (c *NameServerConfig) Build() (*dns.NameServer, error) {
	if c.Address == nil {
		return nil, newError("NameServer address is not specified.")
	}

	var domains []*dns.NameServer_PriorityDomain

	for _, d := range c.Domains {
		newPattern, err := compressPattern(d, externalRules)
		if err != nil {
			return nil, newError("invalid domain rule: ", d).Base(err)
		}
		domains = append(domains, &dns.NameServer_PriorityDomain{
			Type:   dns.DomainMatchingType_New,
			Domain: newPattern,
		})
	}

	geoipList, err := toCidrList(c.ExpectIPs)
	if err != nil {
		return nil, newError("invalid ip rule: ", c.ExpectIPs).Base(err)
	}

	return &dns.NameServer{
		Address: &net.Endpoint{
			Network: net.Network_UDP,
			Address: c.Address.Build(),
			Port:    uint32(c.Port),
		},
		PrioritizedDomain: domains,
		Geoip:             geoipList,
	}, nil
}

// DnsConfig is a JSON serializable object for dns.Config.
type DnsConfig struct {
	Servers  []*NameServerConfig `json:"servers"`
	Hosts    map[string]*Address `json:"hosts"`
	ClientIP *Address            `json:"clientIp"`
	Tag      string              `json:"tag"`
	Fake     *FakeIPConfig       `json:"fake"`
}

func getHostMapping(addr *Address, pattern string) (*dns.Config_HostMapping, error) {
	item := &dns.Config_HostMapping{
		Type: dns.DomainMatchingType_New,
	}
	if addr.Family().IsIP() {
		item.Ip = [][]byte{[]byte(addr.IP())}
	} else {
		item.ProxiedDomain = addr.Domain()
	}
	newPattern, err := compressPattern(pattern, externalRules)
	if err != nil {
		return nil, newError("invalid domain rule: ", pattern).Base(err)
	}
	item.Domain = newPattern
	return item, nil
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

	for pattern, address := range c.Hosts {
		mapping, err := getHostMapping(address, pattern)
		if err != nil {
			return nil, newError("failed to build host rules").Base(err)
		}
		config.StaticHosts = append(config.StaticHosts, mapping)
	}

	if c.Fake != nil {
		config.Fake = new(dns.Config_Fake)
		if c.Fake.FakeNet == "" {
			config.Fake.FakeNet = "224.0.0.0/8"
		} else {
			config.Fake.FakeNet = c.Fake.FakeNet
		}
		if c.Fake.FakeRules != nil {
			fakeRules := make([]string, len(c.Fake.FakeRules))
			i := 0
			for _, pattern := range c.Fake.FakeRules {
				newPattern, err := compressPattern(pattern, externalRules)
				if err == nil {
					fakeRules[i] = newPattern
					i++
				}
			}
			config.Fake.FakeRules = fakeRules[:i]
		}
	}

  if len(externalRules) != 0 {
    config.ExternalRules = make(map[string]*dns.ConfigPatterns)
    for key, value := range externalRules {
      config.ExternalRules[key] = &dns.ConfigPatterns {
        Patterns: value,
      }
    }
  }

	return config, nil
}
