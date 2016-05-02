package hub

import (
	"net"
	"time"
)

type ConnectionHandler func(*Connection)

type Connection struct {
	conn     net.Conn
	listener *TCPHub
}

func (this *Connection) Read(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, ErrorClosedConnection
	}

	return this.conn.Read(b)
}

func (this *Connection) Write(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, ErrorClosedConnection
	}
	return this.conn.Write(b)
}

func (this *Connection) Close() error {
	if this == nil || this.conn == nil {
		return ErrorClosedConnection
	}
	err := this.conn.Close()
	return err
}

func (this *Connection) Release() {
	if this == nil || this.listener == nil {
		return
	}

	this.Close()
	this.conn = nil
	this.listener = nil
}

func (this *Connection) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *Connection) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *Connection) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}

func (this *Connection) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

func (this *Connection) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}
