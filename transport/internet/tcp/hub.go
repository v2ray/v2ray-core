package tcp

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"time"

	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
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
	tlsConfig     *tls.Config
	authConfig    internet.ConnectionAuthenticator
	config        *Config
}

func ListenTCP(address v2net.Address, port v2net.Port, options internet.ListenOptions) (internet.Listener, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	networkSettings, err := options.Stream.GetEffectiveNetworkSettings()
	if err != nil {
		return nil, err
	}
	tcpSettings := networkSettings.(*Config)

	l := &TCPListener{
		acccepting:    true,
		listener:      listener,
		awaitingConns: make(chan *ConnectionWithError, 32),
		config:        tcpSettings,
	}
	if options.Stream != nil && options.Stream.HasSecuritySettings() {
		securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
		if err != nil {
			log.Error("TCP: Failed to get security config: ", err)
			return nil, err
		}
		tlsConfig, ok := securitySettings.(*v2tls.Config)
		if ok {
			l.tlsConfig = tlsConfig.GetTLSConfig()
		}
	}
	if tcpSettings.HeaderSettings != nil {
		headerConfig, err := tcpSettings.HeaderSettings.GetInstance()
		if err != nil {
			return nil, errors.New("TCP: Failed to get header settings: " + err.Error())
		}
		auth, err := internet.CreateConnectionAuthenticator(tcpSettings.HeaderSettings.Type, headerConfig)
		if err != nil {
			return nil, errors.New("TCP: Failed to create header authenticator: " + err.Error())
		}
		l.authConfig = auth
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
			conn := connErr.conn
			return NewConnection(internal.ConnectionId{}, conn, this, this.config), nil
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
		if this.tlsConfig != nil {
			conn = tls.Server(conn, this.tlsConfig)
		}
		if this.authConfig != nil {
			conn = this.authConfig.Server(conn)
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

func (this *TCPListener) Put(id internal.ConnectionId, conn net.Conn) {
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

func ListenRawTCP(address v2net.Address, port v2net.Port, options internet.ListenOptions) (internet.Listener, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	// TODO: handle listen options
	return &RawTCPListener{
		accepting: true,
		listener:  listener,
	}, nil
}

func init() {
	internet.TCPListenFunc = ListenTCP
	internet.RawTCPListenFunc = ListenRawTCP
}
