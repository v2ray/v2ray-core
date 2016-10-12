package dns

import (
	"net"

	"v2ray.com/core/common/log"
)

func (this *Config) GetInternalHosts() map[string]net.IP {
	hosts := make(map[string]net.IP)
	for domain, ipOrDomain := range this.GetHosts() {
		address := ipOrDomain.AsAddress()
		if address.Family().IsDomain() {
			log.Warning("DNS: Ignoring domain address in static hosts: ", address.Domain())
			continue
		}
		hosts[domain] = address.IP()
	}
	return hosts
}
