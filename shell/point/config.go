package point

import (
	"io"

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
	StreamSettings         *internet.StreamConfig
	Protocol               string
	Settings               []byte
	AllowPassiveConnection bool
}

type OutboundConnectionConfig struct {
	Protocol       string
	SendThrough    v2net.Address
	StreamSettings *internet.StreamConfig
	Settings       []byte
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
	StreamSettings         *internet.StreamConfig
	Settings               []byte
	AllowPassiveConnection bool
}

type OutboundDetourConfig struct {
	Protocol       string
	SendThrough    v2net.Address
	StreamSettings *internet.StreamConfig
	Tag            string
	Settings       []byte
}

type Config struct {
	Port            v2net.Port
	LogConfig       *log.Config
	RouterConfig    *router.Config
	DNSConfig       *dns.Config
	InboundConfig   *InboundConnectionConfig
	OutboundConfig  *OutboundConnectionConfig
	InboundDetours  []*InboundDetourConfig
	OutboundDetours []*OutboundDetourConfig
	TransportConfig *transport.Config
}

type ConfigLoader func(input io.Reader) (*Config, error)

var (
	configLoader ConfigLoader
)

func LoadConfig(input io.Reader) (*Config, error) {
	if configLoader == nil {
		return nil, common.ErrBadConfiguration
	}
	return configLoader(input)
}
