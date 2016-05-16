// +build json

package dns

import (
	"encoding/json"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

func (this *Config) UnmarshalJSON(data []byte) error {
	type JsonConfig struct {
		Servers []v2net.Address `json:"servers"`
	}
	jsonConfig := new(JsonConfig)
	if err := json.Unmarshal(data, jsonConfig); err != nil {
		return err
	}
	this.NameServers = make([]v2net.Destination, len(jsonConfig.Servers))
	for idx, server := range jsonConfig.Servers {
		this.NameServers[idx] = v2net.UDPDestination(server, v2net.Port(53))
	}

	return nil
}
