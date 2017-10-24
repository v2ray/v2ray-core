package internet

import (
	"net"
)

type Connection interface {
	net.Conn
}
