package config

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type DetourTag string

type ConnectionConfig interface {
	Protocol() string
	Settings() interface{}
}

type LogConfig interface {
	AccessLog() string
}

type InboundDetourConfig interface {
	Protocol() string
	PortRange() v2net.PortRange
	Settings() interface{}
}

type OutboundDetourConfig interface {
	Protocol() string
	Tag() DetourTag
	Settings() interface{}
}

type PointConfig interface {
	Port() uint16
	LogConfig() LogConfig
	InboundConfig() ConnectionConfig
	OutboundConfig() ConnectionConfig
	InboundDetours() []InboundDetourConfig
}
