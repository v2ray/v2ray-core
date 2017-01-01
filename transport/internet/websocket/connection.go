package websocket

import (
	"io"
	"net"
	"time"

	"v2ray.com/core/common/errors"
)

var (
	ErrInvalidConn = errors.New("Invalid Connection.")
)

type ConnectionManager interface {
	Recycle(string, *wsconn)
}

type Connection struct {
	dest     string
	conn     *wsconn
	listener ConnectionManager
	reusable bool
	config   *Config
}

func NewConnection(dest string, conn *wsconn, manager ConnectionManager, config *Config) *Connection {
	return &Connection{
		dest:     dest,
		conn:     conn,
		listener: manager,
		reusable: config.IsConnectionReuse(),
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
		v.listener.Recycle(v.dest, v.conn)
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
	return v.config.IsConnectionReuse() && v.reusable
}

func (v *Connection) SysFd() (int, error) {
	return getSysFd(v.conn)
}

func getSysFd(conn net.Conn) (int, error) {
	return 0, ErrInvalidConn
}
