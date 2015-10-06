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

type PointConfig interface {
	Port() uint16
	InboundConfig() ConnectionConfig
	OutboundConfig() ConnectionConfig
}
