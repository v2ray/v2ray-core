package conf

import (
	"encoding/json"
	"strings"

	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/router"
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
	Fake     *FakeIPConfig       `json:"fake"`
}

var prefixMapper = map[string]string{
	"domain:":  "d",
	"regexp:":  "r",
	"keyword:": "k",
	"full:":    "f",
	"geosite:": "egeosite.dat:",
	"ext:":     "e",
	"geoip:":   "i",
}
var typeMapper = map[router.Domain_Type]string{
	router.Domain_Full:   "f",
	router.Domain_Domain: "d",
	router.Domain_Plain:  "k",
	router.Domain_Regex:  "r",
}

func compressPattern(pattern string) string {
	for prefix, cmd := range prefixMapper {
		if strings.HasPrefix(pattern, prefix) {
			return cmd + pattern[len(prefix):]
		}
	}
	return "f" + pattern // If no prefix, use full match by default
}

func loadExternalRules(pattern string, c *dns.Config) error {
	cmd := pattern[0]
	if cmd != 'e' {
		return nil
	}
	arg := pattern[1:]
	kv := strings.Split(arg, ":")
	if len(kv) != 2 {
		return newError("invalid external resource: ", arg)
	}
	filename, country := kv[0], kv[1]
	domains, err := loadGeositeWithAttr(filename, country)
	if err != nil {
		return newError("invalid external settings from ", filename, ": ", arg).Base(err)
	}
	externRules := &dns.ConfigPatterns{
		Patterns: make([]string, len(domains)),
	}
	index := 0
	for _, d := range domains {
		externRules.Patterns[index] = typeMapper[d.Type] + d.Value
		index++
	}
	if c.ExternalRules == nil {
		c.ExternalRules = make(map[string]*dns.ConfigPatterns)
	}
	c.ExternalRules[arg] = externRules
	return nil
}

func getHostPattern(addr *Address, pattern string, c *dns.Config) {
	item := new(dns.Config_HostMapping)
	if addr.Family().IsIP() {
		item.Ip = [][]byte{[]byte(addr.IP())}
	} else {
		item.ProxiedDomain = addr.Domain()
	}
	item.Pattern = compressPattern(pattern)
	err := loadExternalRules(item.Pattern, c)
	if err == nil {
		c.HostRules = append(c.HostRules, item)
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

	for pattern, address := range c.Hosts {
		getHostPattern(address, pattern, config)
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
				newPattern := compressPattern(pattern)
				err := loadExternalRules(newPattern, config)
				if err == nil {
					fakeRules[i] = newPattern
					i++
				}
			}
			config.Fake.FakeRules = fakeRules[:i]
		}
	}

	return config, nil
}
