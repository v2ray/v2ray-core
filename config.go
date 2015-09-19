package core

type ConnectionConfig interface {
	Protocol() string
	Content() []byte
}

type PointConfig interface {
	Port() uint16
	InboundConfig() ConnectionConfig
	OutboundConfig() ConnectionConfig
}
