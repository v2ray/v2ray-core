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

// Reuser determines whether a connection can be reused or not.
type Reuser struct {
	// userEnabled indicates connection-reuse enabled by user.
	userEnabled bool
	// appEnabled indicates connection-reuse enabled by app.
	appEnabled bool
}

// ReuseConnection returns a tracker for tracking connection reusability.
func ReuseConnection(reuse bool) *Reuser {
	return &Reuser{
		userEnabled: reuse,
		appEnabled:  reuse,
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

// NewConnection creates a new connection.
func NewConnection(id ConnectionID, conn net.Conn, manager ConnectionRecyler, reuser *Reuser) *Connection {
	return &Connection{
		id:       id,
		conn:     conn,
		listener: manager,
		reuser:   reuser,
	}
}

// Read implements net.Conn.Read().
func (v *Connection) Read(b []byte) (int, error) {
	conn := v.underlyingConn()
	if conn == nil {
		return 0, io.EOF
	}

	return conn.Read(b)
}

// Write implement net.Conn.Write().
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

// LocalAddr implements net.Conn.LocalAddr().
func (v *Connection) LocalAddr() net.Addr {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.LocalAddr()
}

// RemoteAddr implements net.Conn.RemoteAddr().
func (v *Connection) RemoteAddr() net.Addr {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.RemoteAddr()
}

// SetDeadline implements net.Conn.SetDeadline().
func (v *Connection) SetDeadline(t time.Time) error {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.SetDeadline(t)
}

// SetReadDeadline implements net.Conn.SetReadDeadline().
func (v *Connection) SetReadDeadline(t time.Time) error {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.SetReadDeadline(t)
}

// SetWriteDeadline implements net.Conn.SetWriteDeadline().
func (v *Connection) SetWriteDeadline(t time.Time) error {
	conn := v.underlyingConn()
	if conn == nil {
		return nil
	}
	return conn.SetWriteDeadline(t)
}

// SetReusable implements internet.Reusable.SetReusable().
func (v *Connection) SetReusable(reusable bool) {
	if v == nil {
		return
	}
	v.reuser.appEnabled = reusable
}

// Reusable implements internet.Reusable.Reusable().
func (v *Connection) Reusable() bool {
	if v == nil {
		return false
	}
	return v.reuser.userEnabled && v.reuser.appEnabled
}

// SysFd implement internet.SysFd.SysFd().
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
