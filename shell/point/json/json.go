package json

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/v2ray/v2ray-core/app/router"
	routerjson "github.com/v2ray/v2ray-core/app/router/json"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
	"github.com/v2ray/v2ray-core/shell/point"
)

// Config is the config for Point server.
type Config struct {
	PortValue            v2net.Port               `json:"port"` // Port of this Point server.
	LogConfigValue       *LogConfig               `json:"log"`
	RouterConfigValue    *routerjson.RouterConfig `json:"routing"`
	InboundConfigValue   *ConnectionConfig        `json:"inbound"`
	OutboundConfigValue  *ConnectionConfig        `json:"outbound"`
	InboundDetoursValue  []*InboundDetourConfig   `json:"inboundDetour"`
	OutboundDetoursValue []*OutboundDetourConfig  `json:"outboundDetour"`
}

func (config *Config) Port() v2net.Port {
	return config.PortValue
}

func (config *Config) LogConfig() point.LogConfig {
	if config.LogConfigValue == nil {
		return nil
	}
	return config.LogConfigValue
}

func (this *Config) RouterConfig() router.Config {
	if this.RouterConfigValue == nil {
		return nil
	}
	return this.RouterConfigValue
}

func (config *Config) InboundConfig() point.ConnectionConfig {
	if config.InboundConfigValue == nil {
		return nil
	}
	return config.InboundConfigValue
}

func (config *Config) OutboundConfig() point.ConnectionConfig {
	if config.OutboundConfigValue == nil {
		return nil
	}
	return config.OutboundConfigValue
}

func (this *Config) InboundDetours() []point.InboundDetourConfig {
	detours := make([]point.InboundDetourConfig, len(this.InboundDetoursValue))
	for idx, detour := range this.InboundDetoursValue {
		detours[idx] = detour
	}
	return detours
}

func (this *Config) OutboundDetours() []point.OutboundDetourConfig {
	detours := make([]point.OutboundDetourConfig, len(this.OutboundDetoursValue))
	for idx, detour := range this.OutboundDetoursValue {
		detours[idx] = detour
	}
	return detours
}

func LoadConfig(file string) (*Config, error) {
	fixedFile := os.ExpandEnv(file)
	rawConfig, err := ioutil.ReadFile(fixedFile)
	if err != nil {
		log.Error("Failed to read server config file (%s): %v", file, err)
		return nil, err
	}

	jsonConfig := &Config{}
	err = json.Unmarshal(rawConfig, jsonConfig)
	if err != nil {
		log.Error("Failed to load server config: %v", err)
		return nil, err
	}

	jsonConfig.InboundConfigValue.Type = proxyconfig.TypeInbound
	jsonConfig.OutboundConfigValue.Type = proxyconfig.TypeOutbound

	return jsonConfig, err
}
