package tcp

import (
	"io"
	"net"
	"time"

	"v2ray.com/core/transport/internet/internal"
)

type ConnectionManager interface {
	Put(internal.ConnectionId, net.Conn)
}

type RawConnection struct {
	net.TCPConn
}

func (v *RawConnection) Reusable() bool {
	return false
}

func (v *RawConnection) SetReusable(b bool) {}

func (v *RawConnection) SysFd() (int, error) {
	return internal.GetSysFd(&v.TCPConn)
}

type Connection struct {
	id       internal.ConnectionId
	conn     net.Conn
	listener ConnectionManager
	reusable bool
	config   *Config
}

func NewConnection(id internal.ConnectionId, conn net.Conn, manager ConnectionManager, config *Config) *Connection {
	return &Connection{
		id:       id,
		conn:     conn,
		listener: manager,
		reusable: config.ConnectionReuse.IsEnabled(),
		config:   config,
	}
}

func (v *Connection) Read(b []byte) (int, error) {
	if v == nil || v.conn == nil {
		return 0, io.EOF
	}

	return v.conn.Read(b)
}

func (v *Connection) Write(b []byte) (int, error) {
	if v == nil || v.conn == nil {
		return 0, io.ErrClosedPipe
	}
	return v.conn.Write(b)
}

func (v *Connection) Close() error {
	if v == nil || v.conn == nil {
		return io.ErrClosedPipe
	}
	if v.Reusable() {
		v.listener.Put(v.id, v.conn)
		return nil
	}
	err := v.conn.Close()
	v.conn = nil
	return err
}

func (v *Connection) LocalAddr() net.Addr {
	return v.conn.LocalAddr()
}

func (v *Connection) RemoteAddr() net.Addr {
	return v.conn.RemoteAddr()
}

func (v *Connection) SetDeadline(t time.Time) error {
	return v.conn.SetDeadline(t)
}

func (v *Connection) SetReadDeadline(t time.Time) error {
	return v.conn.SetReadDeadline(t)
}

func (v *Connection) SetWriteDeadline(t time.Time) error {
	return v.conn.SetWriteDeadline(t)
}

func (v *Connection) SetReusable(reusable bool) {
	v.reusable = reusable
}

func (v *Connection) Reusable() bool {
	return v.config.ConnectionReuse.IsEnabled() && v.reusable
}

func (v *Connection) SysFd() (int, error) {
	return internal.GetSysFd(v.conn)
}
