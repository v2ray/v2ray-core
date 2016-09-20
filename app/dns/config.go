package dns

import (
	"net"

	"v2ray.com/core/common/log"
)

func (this *Config) GetInternalHosts() map[string]net.IP {
	hosts := make(map[string]net.IP)
	for domain, addressPB := range this.GetHosts() {
		address := addressPB.AsAddress()
		if address.Family().IsDomain() {
			log.Warning("DNS: Ignoring domain address in static hosts: ", address.Domain())
			continue
		}
		hosts[domain] = address.IP()
	}
	return hosts
}
