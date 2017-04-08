package dns

import (
	"net"

	"v2ray.com/core/app/log"
)

func (v *Config) GetInternalHosts() map[string]net.IP {
	hosts := make(map[string]net.IP)
	for domain, ipOrDomain := range v.GetHosts() {
		address := ipOrDomain.AsAddress()
		if address.Family().IsDomain() {
			log.Trace(newError("ignoring domain address in static hosts: ", address.Domain()).AtWarning())
			continue
		}
		hosts[domain] = address.IP()
	}
	return hosts
}
