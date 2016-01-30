package point

import (
	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type ConnectionConfig struct {
	Protocol string
	Settings []byte
}

type LogConfig struct {
	AccessLog string
	ErrorLog  string
	LogLevel  log.LogLevel
}

type DnsConfig struct {
	Enabled  bool
	Settings *dns.CacheConfig
}

const (
	AllocationStrategyAlways   = "always"
	AllocationStrategyRandom   = "random"
	AllocationStrategyExternal = "external"
)

type InboundDetourAllocationConfig struct {
	Strategy    string // Allocation strategy of this inbound detour.
	Concurrency int    // Number of handlers (ports) running in parallel.
	Refresh     int    // Number of minutes before a handler is regenerated.
}

type InboundDetourConfig struct {
	Protocol   string
	PortRange  v2net.PortRange
	Tag        string
	Allocation *InboundDetourAllocationConfig
	Settings   []byte
}

type OutboundDetourConfig struct {
	Protocol string
	Tag      string
	Settings []byte
}

type Config struct {
	Port            v2net.Port
	LogConfig       *LogConfig
	RouterConfig    *router.Config
	InboundConfig   *ConnectionConfig
	OutboundConfig  *ConnectionConfig
	InboundDetours  []*InboundDetourConfig
	OutboundDetours []*OutboundDetourConfig
}

type ConfigLoader func(init string) (*Config, error)

var (
	configLoader ConfigLoader
)

func LoadConfig(init string) (*Config, error) {
	if configLoader == nil {
		return nil, ErrorBadConfiguration
	}
	return configLoader(init)
}
