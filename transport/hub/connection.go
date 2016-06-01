package hub

import (
	"net"
	"time"

	"github.com/v2ray/v2ray-core/transport"
)

type ConnectionHandler func(*Connection)

type ConnectionManager interface {
	Recycle(string, net.Conn)
}

type Connection struct {
	dest     string
	conn     net.Conn
	listener ConnectionManager
	reusable bool
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
	if transport.TCPStreamConfig.ConnectionReuse && this.Reusable() {
		this.listener.Recycle(this.dest, this.conn)
		return nil
	}
	return this.conn.Close()
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

func (this *Connection) SetReusable(reusable bool) {
	this.reusable = reusable
}

func (this *Connection) Reusable() bool {
	return this.reusable
}
