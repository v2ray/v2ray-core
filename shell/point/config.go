package point

import (
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
)

type InboundConnectionConfig struct {
	Port                   v2net.Port
	ListenOn               v2net.Address
	StreamSettings         *internet.StreamSettings
	Protocol               string
	Settings               []byte
	AllowPassiveConnection bool
}

type OutboundConnectionConfig struct {
	Protocol       string
	SendThrough    v2net.Address
	StreamSettings *internet.StreamSettings
	Settings       []byte
}

type LogConfig struct {
	AccessLog string
	ErrorLog  string
	LogLevel  log.LogLevel
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
	Protocol               string
	PortRange              v2net.PortRange
	ListenOn               v2net.Address
	Tag                    string
	Allocation             *InboundDetourAllocationConfig
	StreamSettings         *internet.StreamSettings
	Settings               []byte
	AllowPassiveConnection bool
}

type OutboundDetourConfig struct {
	Protocol       string
	SendThrough    v2net.Address
	StreamSettings *internet.StreamSettings
	Tag            string
	Settings       []byte
}

type Config struct {
	Port            v2net.Port
	LogConfig       *LogConfig
	RouterConfig    *router.Config
	DNSConfig       *dns.Config
	InboundConfig   *InboundConnectionConfig
	OutboundConfig  *OutboundConnectionConfig
	InboundDetours  []*InboundDetourConfig
	OutboundDetours []*OutboundDetourConfig
	TransportConfig *transport.Config
}

type ConfigLoader func(init string) (*Config, error)

var (
	configLoader ConfigLoader
)

func LoadConfig(init string) (*Config, error) {
	if configLoader == nil {
		return nil, common.ErrBadConfiguration
	}
	return configLoader(init)
}
