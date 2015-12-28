package mocks

import (
	"github.com/v2ray/v2ray-core/app/router"
	routertesting "github.com/v2ray/v2ray-core/app/router/testing"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/shell/point"
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
	ErrorLogValue  string
	LogLevelValue  log.LogLevel
}

func (config *LogConfig) AccessLog() string {
	return config.AccessLogValue
}

func (this *LogConfig) ErrorLog() string {
	return this.ErrorLogValue
}

func (this *LogConfig) LogLevel() log.LogLevel {
	return this.LogLevelValue
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

type InboundDetourAllocationConfig struct {
	StrategyValue    string
	ConcurrencyValue int
}

func (this *InboundDetourAllocationConfig) Strategy() string {
	return this.StrategyValue
}

func (this *InboundDetourAllocationConfig) Concurrency() int {
	return this.ConcurrencyValue
}

type InboundDetourConfig struct {
	*ConnectionConfig
	PortRangeValue     *PortRange
	TagValue           string
	AllocationStrategy *InboundDetourAllocationConfig
}

func (this *InboundDetourConfig) Allocation() point.InboundDetourAllocationConfig {
	return this.AllocationStrategy
}

func (this *InboundDetourConfig) Tag() string {
	return this.TagValue
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

type Config struct {
	PortValue            v2net.Port
	LogConfigValue       *LogConfig
	RouterConfigValue    *routertesting.RouterConfig
	InboundConfigValue   *ConnectionConfig
	OutboundConfigValue  *ConnectionConfig
	InboundDetoursValue  []*InboundDetourConfig
	OutboundDetoursValue []*OutboundDetourConfig
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

func (this *Config) InboundConfig() point.ConnectionConfig {
	if this.InboundConfigValue == nil {
		return nil
	}
	return this.InboundConfigValue
}

func (this *Config) OutboundConfig() point.ConnectionConfig {
	if this.OutboundConfigValue == nil {
		return nil
	}
	return this.OutboundConfigValue
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
