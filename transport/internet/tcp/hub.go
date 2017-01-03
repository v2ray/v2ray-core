package tcp

import (
	"crypto/tls"
	"net"
	"sync"
	"time"

	"v2ray.com/core/common/errors"
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
			return nil, errors.Base(err).Message("Internet|TCP: Invalid header settings.")
		}
		auth, err := internet.CreateConnectionAuthenticator(tcpSettings.HeaderSettings.Type, headerConfig)
		if err != nil {
			return nil, errors.Base(err).Message("Internet|TCP: Invalid header settings.")
		}
		l.authConfig = auth
	}
	go l.KeepAccepting()
	return l, nil
}

func (v *TCPListener) Accept() (internet.Connection, error) {
	for v.acccepting {
		select {
		case connErr, open := <-v.awaitingConns:
			if !open {
				return nil, ErrClosedListener
			}
			if connErr.err != nil {
				return nil, connErr.err
			}
			conn := connErr.conn
			return NewConnection(internal.ConnectionID{}, conn, v, v.config), nil
		case <-time.After(time.Second * 2):
		}
	}
	return nil, ErrClosedListener
}

func (v *TCPListener) KeepAccepting() {
	for v.acccepting {
		conn, err := v.listener.Accept()
		v.Lock()
		if !v.acccepting {
			v.Unlock()
			break
		}
		if v.tlsConfig != nil {
			conn = tls.Server(conn, v.tlsConfig)
		}
		if v.authConfig != nil {
			conn = v.authConfig.Server(conn)
		}
		select {
		case v.awaitingConns <- &ConnectionWithError{
			conn: conn,
			err:  err,
		}:
		default:
			if conn != nil {
				conn.Close()
			}
		}

		v.Unlock()
	}
}

func (v *TCPListener) Put(id internal.ConnectionID, conn net.Conn) {
	v.Lock()
	defer v.Unlock()
	if !v.acccepting {
		return
	}
	select {
	case v.awaitingConns <- &ConnectionWithError{conn: conn}:
	default:
		conn.Close()
	}
}

func (v *TCPListener) Addr() net.Addr {
	return v.listener.Addr()
}

func (v *TCPListener) Close() error {
	v.Lock()
	defer v.Unlock()
	v.acccepting = false
	v.listener.Close()
	close(v.awaitingConns)
	for connErr := range v.awaitingConns {
		if connErr.conn != nil {
			go connErr.conn.Close()
		}
	}
	return nil
}

func init() {
	internet.TCPListenFunc = ListenTCP
}
