package config

type RouterConfig interface {
	Strategy() string
	Settings() interface{}
}

type ConnectionTag string

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

type PointConfig interface {
	Port() uint16
	LogConfig() LogConfig
	InboundConfig() ConnectionConfig
	OutboundConfig() ConnectionConfig
	InboundDetours() []InboundDetourConfig
}
