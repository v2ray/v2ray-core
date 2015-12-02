package mocks

import (
	routerconfig "github.com/v2ray/v2ray-core/app/router/config"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/shell/point/config"
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
	FromValue v2net.Port
	ToValue   v2net.Port
}

func (this *PortRange) From() v2net.Port {
	return this.FromValue
}

func (this *PortRange) To() v2net.Port {
	return this.ToValue
}

type InboundDetourConfig struct {
	ConnectionConfig
	PortRangeValue *PortRange
}

func (this *InboundDetourConfig) PortRange() v2net.PortRange {
	return this.PortRangeValue
}

type OutboundDetourConfig struct {
	ConnectionConfig
	TagValue string
}

func (this *OutboundDetourConfig) Tag() string {
	return this.TagValue
}

func (config *LogConfig) AccessLog() string {
	return config.AccessLogValue
}

type Config struct {
	PortValue            uint16
	LogConfigValue       *LogConfig
	RouterConfigValue    routerconfig.RouterConfig
	InboundConfigValue   *ConnectionConfig
	OutboundConfigValue  *ConnectionConfig
	InboundDetoursValue  []*InboundDetourConfig
	OutboundDetoursValue []*OutboundDetourConfig
}

func (config *Config) Port() uint16 {
	return config.PortValue
}

func (config *Config) LogConfig() config.LogConfig {
	return config.LogConfigValue
}

func (this *Config) RouterConfig() routerconfig.RouterConfig {
	return this.RouterConfigValue
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

func (this *Config) OutboundDetours() []config.OutboundDetourConfig {
	detours := make([]config.OutboundDetourConfig, len(this.OutboundDetoursValue))
	for idx, detour := range this.OutboundDetoursValue {
		detours[idx] = detour
	}
	return detours
}
