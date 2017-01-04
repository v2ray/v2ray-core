package internal

import (
	"io"
	"net"
	"sync"
	"time"

	v2net "v2ray.com/core/common/net"
)

// ConnectionID is the ID of a connection.
type ConnectionID struct {
	Local      v2net.Address
	Remote     v2net.Address
	RemotePort v2net.Port
}

// NewConnectionID creates a new ConnectionId.
func NewConnectionID(source v2net.Address, dest v2net.Destination) ConnectionID {
	return ConnectionID{
		Local:      source,
		Remote:     dest.Address,
		RemotePort: dest.Port,
	}
}

type Reuser struct {
	userEnabled bool
	appEnable   bool
}

func ReuseConnection(reuse bool) *Reuser {
	return &Reuser{
		userEnabled: reuse,
		appEnable:   reuse,
	}
}

// Connection is an implementation of net.Conn with re-usability.
type Connection struct {
	sync.RWMutex
	id       ConnectionID
	conn     net.Conn
	listener ConnectionRecyler
	reuser   *Reuser
}

func NewConnection(id ConnectionID, conn net.Conn, manager ConnectionRecyler, reuser *Reuser) *Connection {
	return &Connection{
		id:       id,
		conn:     conn,
		listener: manager,
		reuser:   reuser,
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

// Close implements net.Conn.Close(). If the connection is reusable, the underlying connection will be recycled.
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
	v.reuser.appEnable = reusable
}

func (v *Connection) Reusable() bool {
	if v == nil {
		return false
	}
	return v.reuser.userEnabled && v.reuser.appEnable
}

func (v *Connection) SysFd() (int, error) {
	conn := v.underlyingConn()
	if conn == nil {
		return 0, io.ErrClosedPipe
	}
	return GetSysFd(conn)
}

func (v *Connection) underlyingConn() net.Conn {
	if v == nil {
		return nil
	}

	v.RLock()
	defer v.RUnlock()

	return v.conn
}
