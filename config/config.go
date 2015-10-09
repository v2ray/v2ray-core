package config

type Type string

const (
	TypeInbound  = Type("inbound")
	TypeOutbound = Type("outbound")
)

type ConnectionConfig interface {
	Protocol() string
	Settings(configType Type) interface{}
}

type LogConfig interface {
	AccessLog() string
}

type PointConfig interface {
	Port() uint16
	LogConfig() LogConfig
	InboundConfig() ConnectionConfig
	OutboundConfig() ConnectionConfig
}
