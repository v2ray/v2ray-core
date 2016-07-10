package tls

import (
	"crypto/tls"
)

type Connection struct {
	*tls.Conn
}

func (this *Connection) Reusable() bool {
	return false
}

func (this *Connection) SetReusable(bool) {}

func NewConnection(conn *tls.Conn) *Connection {
	return &Connection{
		Conn: conn,
	}
}
