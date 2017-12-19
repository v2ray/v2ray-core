package dns

import (
	"net"
)

func (c *Config) GetInternalHosts() map[string]net.IP {
	hosts := make(map[string]net.IP)
	for domain, ipOrDomain := range c.GetHosts() {
		address := ipOrDomain.AsAddress()
		if address.Family().IsDomain() {
			newError("ignoring domain address in static hosts: ", address.Domain()).AtWarning().WriteToLog()
			continue
		}
		hosts[domain] = address.IP()
	}
	return hosts
}
