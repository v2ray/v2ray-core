package json

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/log"
)

type ConnectionConfig struct {
	ProtocolString string `json:"protocol"`
	File           string `json:"file"`
}

func (config *ConnectionConfig) Protocol() string {
	return config.ProtocolString
}

func (config *ConnectionConfig) Content() []byte {
	if len(config.File) == 0 {
		return nil
	}
	content, err := ioutil.ReadFile(config.File)
	if err != nil {
		panic(log.Error("Failed to read config file (%s): %v", config.File, err))
	}
	return content
}

// Config is the config for Point server.
type Config struct {
	PortValue           uint16            `json:"port"` // Port of this Point server.
	InboundConfigValue  *ConnectionConfig `json:"inbound"`
	OutboundConfigValue *ConnectionConfig `json:"outbound"`
}

func (config *Config) Port() uint16 {
	return config.PortValue
}

func (config *Config) InboundConfig() core.ConnectionConfig {
	return config.InboundConfigValue
}

func (config *Config) OutboundConfig() core.ConnectionConfig {
	return config.OutboundConfigValue
}

func LoadConfig(file string) (*Config, error) {
	fixedFile := os.ExpandEnv(file)
	rawConfig, err := ioutil.ReadFile(fixedFile)
	if err != nil {
		log.Error("Failed to read point config file (%s): %v", file, err)
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(rawConfig, config)

	if !filepath.IsAbs(config.InboundConfigValue.File) && len(config.InboundConfigValue.File) > 0 {
		config.InboundConfigValue.File = filepath.Join(filepath.Dir(fixedFile), config.InboundConfigValue.File)
	}

	if !filepath.IsAbs(config.OutboundConfigValue.File) && len(config.OutboundConfigValue.File) > 0 {
		config.OutboundConfigValue.File = filepath.Join(filepath.Dir(fixedFile), config.OutboundConfigValue.File)
	}

	return config, err
}
