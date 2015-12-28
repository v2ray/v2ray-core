package point

import (
	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type ConnectionConfig interface {
	Protocol() string
	Settings() interface{}
}

type LogConfig interface {
	AccessLog() string
	ErrorLog() string
	LogLevel() log.LogLevel
}

type DnsConfig interface {
	Enabled() bool
	Settings() dns.CacheConfig
}

const (
	AllocationStrategyAlways   = "always"
	AllocationStrategyRandom   = "random"
	AllocationStrategyExternal = "external"
)

type InboundDetourAllocationConfig interface {
	Strategy() string
	Concurrency() int
}

type InboundDetourConfig interface {
	Protocol() string
	PortRange() v2net.PortRange
	Tag() string
	Allocation() InboundDetourAllocationConfig
	Settings() interface{}
}

type OutboundDetourConfig interface {
	Protocol() string
	Tag() string
	Settings() interface{}
}

type PointConfig interface {
	Port() v2net.Port
	LogConfig() LogConfig
	RouterConfig() router.Config
	InboundConfig() ConnectionConfig
	OutboundConfig() ConnectionConfig
	InboundDetours() []InboundDetourConfig
	OutboundDetours() []OutboundDetourConfig
}
