package mocks

import (
	"github.com/v2ray/v2ray-core/app/point/config"
)

type ConnectionConfig struct {
	ProtocolValue string
	SettingsValue interface{}
}

func (config *ConnectionConfig) Protocol() string {
	return config.ProtocolValue
}

func (config *ConnectionConfig) Settings() interface{} {
	return config.SettingsValue
}

type LogConfig struct {
	AccessLogValue string
}

func (config *LogConfig) AccessLog() string {
	return config.AccessLogValue
}

type Config struct {
	PortValue           uint16
	LogConfigValue      *LogConfig
	InboundConfigValue  *ConnectionConfig
	OutboundConfigValue *ConnectionConfig
}

func (config *Config) Port() uint16 {
	return config.PortValue
}

func (config *Config) LogConfig() config.LogConfig {
	return config.LogConfigValue
}

func (config *Config) InboundConfig() config.ConnectionConfig {
	return config.InboundConfigValue
}

func (config *Config) OutboundConfig() config.ConnectionConfig {
	return config.OutboundConfigValue
}
