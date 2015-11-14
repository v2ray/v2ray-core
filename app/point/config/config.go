package config

type DetourTag string

type ConnectionConfig interface {
	Protocol() string
	Settings() interface{}
}

type LogConfig interface {
	AccessLog() string
}

type PortRange interface {
	From() uint16
	To() uint16
}

type InboundDetourConfig interface {
	Protocol() string
	PortRange() PortRange
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
