package hub

import (
	"net"
	"time"

	"github.com/v2ray/v2ray-core/common"
)

type ConnectionHandler func(Connection)

type Connection interface {
	common.Releasable

	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	CloseRead() error
	CloseWrite() error
}

type TCPConnection struct {
	conn     *net.TCPConn
	listener *TCPHub
}

func (this *TCPConnection) Read(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, ErrorClosedConnection
	}

	return this.conn.Read(b)
}

func (this *TCPConnection) Write(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, ErrorClosedConnection
	}
	return this.conn.Write(b)
}

func (this *TCPConnection) Close() error {
	if this == nil || this.conn == nil {
		return ErrorClosedConnection
	}
	err := this.conn.Close()
	return err
}

func (this *TCPConnection) Release() {
	if this == nil || this.listener == nil {
		return
	}

	this.Close()
	this.conn = nil
	this.listener = nil
}

func (this *TCPConnection) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *TCPConnection) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *TCPConnection) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}

func (this *TCPConnection) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

func (this *TCPConnection) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}

func (this *TCPConnection) CloseRead() error {
	if this == nil || this.conn == nil {
		return nil
	}
	return this.conn.CloseRead()
}

func (this *TCPConnection) CloseWrite() error {
	if this == nil || this.conn == nil {
		return nil
	}
	return this.conn.CloseWrite()
}
