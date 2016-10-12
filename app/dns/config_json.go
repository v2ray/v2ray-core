// +build json

package dns

import (
	"encoding/json"

	v2net "v2ray.com/core/common/net"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Servers []*v2net.IPOrDomain          `json:"servers"`
		Hosts   map[string]*v2net.IPOrDomain `json:"hosts"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.NameServers = make([]*v2net.Endpoint, len(jsonConfig.Servers))
	for idx, server := range jsonConfig.Servers {
		this.NameServers[idx] = &v2net.Endpoint{
			Network: v2net.Network_UDP,
			Address: server,
			Port:    53,
		}
	}

	if jsonConfig.Hosts != nil {
		this.Hosts = jsonConfig.Hosts
	}

	return nil
}
