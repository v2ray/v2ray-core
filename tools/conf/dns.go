package conf

import (
	"v2ray.com/core/app/dns"
	v2net "v2ray.com/core/common/net"
)

type DnsConfig struct {
	Servers []*Address          `json:"servers"`
	Hosts   map[string]*Address `json:"hosts"`
}

func (v *DnsConfig) Build() *dns.Config {
	config := new(dns.Config)
	config.NameServers = make([]*v2net.Endpoint, len(v.Servers))
	for idx, server := range v.Servers {
		config.NameServers[idx] = &v2net.Endpoint{
			Network: v2net.Network_UDP,
			Address: server.Build(),
			Port:    53,
		}
	}

	if v.Hosts != nil {
		config.Hosts = make(map[string]*v2net.IPOrDomain)
		for domain, ip := range v.Hosts {
			config.Hosts[domain] = ip.Build()
		}
	}

	return config
}
