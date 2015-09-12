package vmess

import (
	"encoding/json"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/net"
)

type VMessInboundConfig struct {
	AllowedClients []core.VUser `json:"clients"`
}

func loadInboundConfig(rawConfig []byte) (VMessInboundConfig, error) {
	config := VMessInboundConfig{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}

type VNextConfig struct {
	Address string       `json:"address"`
	Port    uint16       `json:"port"`
	Users   []core.VUser `json:"users"`
}

func (config VNextConfig) ToVNextServer() VNextServer {
	return VNextServer{
		v2net.DomainAddress(config.Address, config.Port),
		config.Users}
}

type VMessOutboundConfig struct {
	VNextList []VNextConfig
}

func loadOutboundConfig(rawConfig []byte) (VMessOutboundConfig, error) {
	config := VMessOutboundConfig{}
	err := json.Unmarshal(rawConfig, &config)
	return config, err
}
