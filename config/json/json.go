package json

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/config"
)

type ConnectionConfig struct {
	ProtocolString  string          `json:"protocol"`
	SettingsMessage json.RawMessage `json:"settings"`
}

func (config *ConnectionConfig) Protocol() string {
	return config.ProtocolString
}

func (config *ConnectionConfig) Settings(configType config.Type) interface{} {
	creator, found := configCache[getConfigKey(config.Protocol(), configType)]
	if !found {
		panic("Unknown protocol " + config.Protocol())
	}
	configObj := creator()
	err := json.Unmarshal(config.SettingsMessage, configObj)
	if err != nil {
		log.Error("Unable to parse connection config: %v", err)
		panic("Failed to parse connection config.")
	}
	return configObj
}

type LogConfig struct {
	AccessLogValue string `json:"access"`
}

func (config *LogConfig) AccessLog() string {
	return config.AccessLogValue
}

// Config is the config for Point server.
type Config struct {
	PortValue           uint16            `json:"port"` // Port of this Point server.
	LogConfigValue      *LogConfig        `json:"log"`
	InboundConfigValue  *ConnectionConfig `json:"inbound"`
	OutboundConfigValue *ConnectionConfig `json:"outbound"`
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

func LoadConfig(file string) (*Config, error) {
	fixedFile := os.ExpandEnv(file)
	rawConfig, err := ioutil.ReadFile(fixedFile)
	if err != nil {
		log.Error("Failed to read server config file (%s): %v", file, err)
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(rawConfig, config)
	if err != nil {
		log.Error("Failed to load server config: %v", err)
		return nil, err
	}

	return config, err
}
