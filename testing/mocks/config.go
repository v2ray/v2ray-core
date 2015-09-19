package mocks

import (
	"github.com/v2ray/v2ray-core"
)

type ConnectionConfig struct {
	ProtocolValue string
	ContentValue  []byte
}

func (config *ConnectionConfig) Protocol() string {
	return config.ProtocolValue
}

func (config *ConnectionConfig) Content() []byte {
	return config.ContentValue
}

type Config struct {
	PortValue           uint16
	InboundConfigValue  *ConnectionConfig
	OutboundConfigValue *ConnectionConfig
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
