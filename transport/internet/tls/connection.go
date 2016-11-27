package tls

import (
	"crypto/tls"
)

type Connection struct {
	*tls.Conn
}

func (v *Connection) Reusable() bool {
	return false
}

func (v *Connection) SetReusable(bool) {}

func NewConnection(conn *tls.Conn) *Connection {
	return &Connection{
		Conn: conn,
	}
}
