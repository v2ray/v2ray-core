package tcp

import (
	"io"
	"net"
	"sync"
	"time"

	"v2ray.com/core/transport/internet/internal"
)

type Connection struct {
	sync.RWMutex
	id       internal.ConnectionID
	reusable bool
	conn     net.Conn
	listener internal.ConnectionRecyler
	config   *Config
}

func NewConnection(id internal.ConnectionID, conn net.Conn, manager internal.ConnectionRecyler, config *Config) *Connection {
	return &Connection{
		id:       id,
		conn:     conn,
		listener: manager,
		reusable: config.IsConnectionReuse(),
		config:   config,
	}
}

func (v *Connection) Read(b []byte) (int, error) {
	conn := v.underlyingConn()
	if conn == nil {
		return 0, io.EOF
	}

	return conn.Read(b)
}

func (v *Connection) Write(b []byte) (int, error) {
	conn := v.underlyingConn()
	if conn == nil {
		return 0, io.ErrClosedPipe
	}
	return conn.Write(b)
}

func (v *Connection) Close() error {
	if v == nil {
		return io.ErrClosedPipe
	}

	v.Lock()
	defer v.Unlock()
	if v.conn == nil {
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
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.LocalAddr()
}

func (v *Connection) RemoteAddr() net.Addr {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.RemoteAddr()
}

func (v *Connection) SetDeadline(t time.Time) error {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.SetDeadline(t)
}

func (v *Connection) SetReadDeadline(t time.Time) error {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.SetReadDeadline(t)
}

func (v *Connection) SetWriteDeadline(t time.Time) error {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.SetWriteDeadline(t)
}

func (v *Connection) SetReusable(reusable bool) {
	if v == nil {
		return
	}
	v.reusable = reusable
}

func (v *Connection) Reusable() bool {
	if v == nil {
		return false
	}
	return v.config.IsConnectionReuse() && v.reusable
}

func (v *Connection) SysFd() (int, error) {
	conn := v.underlyingConn()
	if conn == nil {
		return 0, io.ErrClosedPipe
	}
	return internal.GetSysFd(conn)
}

func (v *Connection) underlyingConn() net.Conn {
	if v == nil {
		return nil
	}

	v.RLock()
	defer v.RUnlock()

	return v.conn
}
