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

type InboundDetour interface {
	Protocol() string
	Port()
}

type PointConfig interface {
	Port() uint16
	LogConfig() LogConfig
	InboundConfig() ConnectionConfig
	OutboundConfig() ConnectionConfig
}
