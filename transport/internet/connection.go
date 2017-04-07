package internet

import (
	"net"
)

type ConnectionHandler func(Connection)

type Connection interface {
	net.Conn
}

type SysFd interface {
	SysFd() (int, error)
}
