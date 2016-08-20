package tcp

import (
	"errors"
	"net"
	"sync"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

var (
	ErrClosedListener = errors.New("Listener is closed.")
)

type ConnectionWithError struct {
	conn net.Conn
	err  error
}

type TCPListener struct {
	sync.Mutex
	acccepting    bool
	listener      *net.TCPListener
	awaitingConns chan *ConnectionWithError
}

func ListenTCP(address v2net.Address, port v2net.Port) (internet.Listener, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	l := &TCPListener{
		acccepting:    true,
		listener:      listener,
		awaitingConns: make(chan *ConnectionWithError, 32),
	}
	go l.KeepAccepting()
	return l, nil
}

func (this *TCPListener) Accept() (internet.Connection, error) {
	for this.acccepting {
		select {
		case connErr, open := <-this.awaitingConns:
			if !open {
				return nil, ErrClosedListener
			}
			if connErr.err != nil {
				return nil, connErr.err
			}
			return NewConnection("", connErr.conn, this), nil
		case <-time.After(time.Second * 2):
		}
	}
	return nil, ErrClosedListener
}

func (this *TCPListener) KeepAccepting() {
	for this.acccepting {
		conn, err := this.listener.Accept()
		this.Lock()
		if !this.acccepting {
			this.Unlock()
			break
		}
		select {
		case this.awaitingConns <- &ConnectionWithError{
			conn: conn,
			err:  err,
		}:
		default:
			if conn != nil {
				conn.Close()
			}
		}

		this.Unlock()
	}
}

func (this *TCPListener) Recycle(dest string, conn net.Conn) {
	this.Lock()
	defer this.Unlock()
	if !this.acccepting {
		return
	}
	select {
	case this.awaitingConns <- &ConnectionWithError{conn: conn}:
	default:
		conn.Close()
	}
}

func (this *TCPListener) Addr() net.Addr {
	return this.listener.Addr()
}

func (this *TCPListener) Close() error {
	this.Lock()
	defer this.Unlock()
	this.acccepting = false
	this.listener.Close()
	close(this.awaitingConns)
	for connErr := range this.awaitingConns {
		if connErr.conn != nil {
			go connErr.conn.Close()
		}
	}
	return nil
}

type RawTCPListener struct {
	accepting bool
	listener  *net.TCPListener
}

func (this *RawTCPListener) Accept() (internet.Connection, error) {
	conn, err := this.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return &RawConnection{
		TCPConn: *conn,
	}, nil
}

func (this *RawTCPListener) Addr() net.Addr {
	return this.listener.Addr()
}

func (this *RawTCPListener) Close() error {
	this.accepting = false
	this.listener.Close()
	return nil
}

func ListenRawTCP(address v2net.Address, port v2net.Port) (internet.Listener, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	return &RawTCPListener{
		accepting: true,
		listener:  listener,
	}, nil
}

func init() {
	internet.TCPListenFunc = ListenTCP
	internet.RawTCPListenFunc = ListenRawTCP
}
