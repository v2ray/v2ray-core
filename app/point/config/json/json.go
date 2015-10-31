package json

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/v2ray/v2ray-core/app/point/config"
	"github.com/v2ray/v2ray-core/common/log"
	proxyconfig "github.com/v2ray/v2ray-core/proxy/common/config"
)

// Config is the config for Point server.
type Config struct {
	PortValue           uint16                 `json:"port"` // Port of this Point server.
	LogConfigValue      *LogConfig             `json:"log"`
	InboundConfigValue  *ConnectionConfig      `json:"inbound"`
	OutboundConfigValue *ConnectionConfig      `json:"outbound"`
	InboundDetoursValue []*InboundDetourConfig `json:"inboundDetour"`
}

func (config *Config) Port() uint16 {
	return config.PortValue
}

func (config *Config) LogConfig() config.LogConfig {
	if config.LogConfigValue == nil {
		return nil
	}
	return config.LogConfigValue
}

func (config *Config) InboundConfig() config.ConnectionConfig {
	if config.InboundConfigValue == nil {
		return nil
	}
	return config.InboundConfigValue
}

func (config *Config) OutboundConfig() config.ConnectionConfig {
	if config.OutboundConfigValue == nil {
		return nil
	}
	return config.OutboundConfigValue
}

func (this *Config) InboundDetours() []config.InboundDetourConfig {
	detours := make([]config.InboundDetourConfig, len(this.InboundDetoursValue))
	for idx, detour := range this.InboundDetoursValue {
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
