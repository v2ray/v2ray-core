package mocks

import (
	routerconfig "github.com/v2ray/v2ray-core/app/router/config"
	routertestingconfig "github.com/v2ray/v2ray-core/app/router/config/testing"
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
	*ConnectionConfig
	PortRangeValue *PortRange
}

func (this *InboundDetourConfig) PortRange() v2net.PortRange {
	return this.PortRangeValue
}

type OutboundDetourConfig struct {
	*ConnectionConfig
	TagValue string
}

func (this *OutboundDetourConfig) Tag() string {
	return this.TagValue
}

func (config *LogConfig) AccessLog() string {
	return config.AccessLogValue
}

type Config struct {
	PortValue            v2net.Port
	LogConfigValue       *LogConfig
	RouterConfigValue    *routertestingconfig.RouterConfig
	InboundConfigValue   *ConnectionConfig
	OutboundConfigValue  *ConnectionConfig
	InboundDetoursValue  []*InboundDetourConfig
	OutboundDetoursValue []*OutboundDetourConfig
}

func (config *Config) Port() v2net.Port {
	return config.PortValue
}

func (config *Config) LogConfig() config.LogConfig {
	if config.LogConfigValue == nil {
		return nil
	}
	return config.LogConfigValue
}

func (this *Config) RouterConfig() routerconfig.RouterConfig {
	if this.RouterConfigValue == nil {
		return nil
	}
	return this.RouterConfigValue
}

func (this *Config) InboundConfig() config.ConnectionConfig {
	if this.InboundConfigValue == nil {
		return nil
	}
	return this.InboundConfigValue
}

func (this *Config) OutboundConfig() config.ConnectionConfig {
	if this.OutboundConfigValue == nil {
		return nil
	}
	return this.OutboundConfigValue
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
