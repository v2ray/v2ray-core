package tcp

import (
	"io"
	"net"
	"time"

	"v2ray.com/core/transport/internet/internal"
)

type ConnectionManager interface {
	Recycle(string, net.Conn)
}

type RawConnection struct {
	net.TCPConn
}

func (this *RawConnection) Reusable() bool {
	return false
}

func (this *RawConnection) SetReusable(b bool) {}

func (this *RawConnection) SysFd() (int, error) {
	return internal.GetSysFd(&this.TCPConn)
}

type Connection struct {
	dest     string
	conn     net.Conn
	listener ConnectionManager
	reusable bool
	config   *Config
}

func NewConnection(dest string, conn net.Conn, manager ConnectionManager, config *Config) *Connection {
	return &Connection{
		dest:     dest,
		conn:     conn,
		listener: manager,
		reusable: config.ConnectionReuse.IsEnabled(),
		config:   config,
	}
}

func (this *Connection) Read(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, io.EOF
	}

	return this.conn.Read(b)
}

func (this *Connection) Write(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, io.ErrClosedPipe
	}
	return this.conn.Write(b)
}

func (this *Connection) Close() error {
	if this == nil || this.conn == nil {
		return io.ErrClosedPipe
	}
	if this.Reusable() {
		this.listener.Recycle(this.dest, this.conn)
		return nil
	}
	err := this.conn.Close()
	this.conn = nil
	return err
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
	if !this.config.ConnectionReuse.IsEnabled() {
		return
	}
	this.reusable = reusable
}

func (this *Connection) Reusable() bool {
	return this.reusable
}

func (this *Connection) SysFd() (int, error) {
	return internal.GetSysFd(this.conn)
}
