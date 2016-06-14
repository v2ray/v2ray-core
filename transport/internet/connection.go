package internet

import (
	"net"
)

type ConnectionHandler func(Connection)

type Reusable interface {
	Reusable() bool
	SetReusable(reuse bool)
}

type StreamConnectionType int

var (
	StreamConnectionTypeRawTCP StreamConnectionType = 1
	StreamConnectionTypeTCP    StreamConnectionType = 2
	StreamConnectionTypeKCP    StreamConnectionType = 4
)

type StreamSettings struct {
	Type StreamConnectionType
}

func (this *StreamSettings) IsCapableOf(streamType StreamConnectionType) bool {
	return (this.Type & streamType) == streamType
}

type Connection interface {
	net.Conn
	Reusable
}
