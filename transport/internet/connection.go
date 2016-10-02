package internet

import (
	"net"
)

type ConnectionHandler func(Connection)

type Reusable interface {
	Reusable() bool
	SetReusable(reuse bool)
}

type Connection interface {
	net.Conn
	Reusable
}

type SysFd interface {
	SysFd() (int, error)
}
