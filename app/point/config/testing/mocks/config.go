package mocks

import (
	"github.com/v2ray/v2ray-core/app/point/config"
	v2net "github.com/v2ray/v2ray-core/common/net"
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

type PortRange struct {
	FromValue uint16
	ToValue   uint16
}

func (this *PortRange) From() uint16 {
	return this.FromValue
}

func (this *PortRange) To() uint16 {
	return this.ToValue
}

type InboundDetourConfig struct {
	ConnectionConfig
	PortRangeValue *PortRange
}

func (this *InboundDetourConfig) PortRange() v2net.PortRange {
	return this.PortRangeValue
}

func (config *LogConfig) AccessLog() string {
	return config.AccessLogValue
}

type Config struct {
	PortValue           uint16
	LogConfigValue      *LogConfig
	InboundConfigValue  *ConnectionConfig
	OutboundConfigValue *ConnectionConfig
	InboundDetoursValue []*InboundDetourConfig
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

func (this *Config) InboundDetours() []config.InboundDetourConfig {
	detours := make([]config.InboundDetourConfig, len(this.InboundDetoursValue))
	for idx, detour := range this.InboundDetoursValue {
		detours[idx] = detour
	}
	return detours
}
