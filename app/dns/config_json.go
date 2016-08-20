// +build json

package dns

import (
	"encoding/json"
	"errors"
	"net"

	v2net "v2ray.com/core/common/net"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Servers []v2net.AddressJson          `json:"servers"`
		Hosts   map[string]v2net.AddressJson `json:"hosts"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.NameServers = make([]v2net.Destination, len(jsonConfig.Servers))
	for idx, server := range jsonConfig.Servers {
		this.NameServers[idx] = v2net.UDPDestination(server.Address, v2net.Port(53))
	}

	if jsonConfig.Hosts != nil {
		this.Hosts = make(map[string]net.IP)
		for domain, ip := range jsonConfig.Hosts {
			if ip.Address.Family().IsDomain() {
				return errors.New(ip.Address.String() + " is not an IP.")
			}
			this.Hosts[domain] = ip.Address.IP()
		}
	}

	return nil
}
