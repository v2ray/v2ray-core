package hub

import (
	"errors"
	"net"
	"time"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

var (
	ErrorClosedConnection = errors.New("Connection already closed.")
)

type TCPConn struct {
	conn     *net.TCPConn
	listener *TCPHub
}

func (this *TCPConn) Read(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, ErrorClosedConnection
	}
	return this.conn.Read(b)
}

func (this *TCPConn) Write(b []byte) (int, error) {
	if this == nil || this.conn == nil {
		return 0, ErrorClosedConnection
	}
	return this.conn.Write(b)
}

func (this *TCPConn) Close() error {
	if this == nil || this.conn == nil {
		return ErrorClosedConnection
	}
	err := this.conn.Close()
	return err
}

func (this *TCPConn) Release() {
	if this == nil || this.listener == nil {
		return
	}

	this.Close()
	this.conn = nil
	this.listener = nil
}

func (this *TCPConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *TCPConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *TCPConn) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}

func (this *TCPConn) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

func (this *TCPConn) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}

func (this *TCPConn) CloseRead() error {
	if this == nil || this.conn == nil {
		return nil
	}
	return this.conn.CloseRead()
}

func (this *TCPConn) CloseWrite() error {
	if this == nil || this.conn == nil {
		return nil
	}
	return this.conn.CloseWrite()
}

type TCPHub struct {
	listener     *net.TCPListener
	connCallback func(*TCPConn)
	accepting    bool
}

func ListenTCP(port v2net.Port, callback func(*TCPConn)) (*TCPHub, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
	tcpListener := &TCPHub{
		listener:     listener,
		connCallback: callback,
	}
	go tcpListener.start()
	return tcpListener, nil
}

func (this *TCPHub) Close() {
	this.accepting = false
	this.listener.Close()
	this.listener = nil
}

func (this *TCPHub) start() {
	this.accepting = true
	for this.accepting {
		conn, err := this.listener.AcceptTCP()
		if err != nil {
			if this.accepting {
				log.Warning("Listener: Failed to accept new TCP connection: ", err)
			}
			continue
		}
		go this.connCallback(&TCPConn{
			conn:     conn,
			listener: this,
		})
	}
}

func (this *TCPHub) recycle(conn *net.TCPConn) {

}
