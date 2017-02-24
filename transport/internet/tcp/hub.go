package tcp

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
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

func ListenTCP(ctx context.Context, address v2net.Address, port v2net.Port) (internet.Listener, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	log.Info("TCP|Listener: Listening on ", address, ":", port)
	networkSettings := internet.TransportSettingsFromContext(ctx)
	tcpSettings := networkSettings.(*Config)

	l := &TCPListener{
		acccepting:    true,
		listener:      listener,
		awaitingConns: make(chan *ConnectionWithError, 32),
		config:        tcpSettings,
	}
	if securitySettings := internet.SecuritySettingsFromContext(ctx); securitySettings != nil {
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
		auth, err := internet.CreateConnectionAuthenticator(headerConfig)
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
			return internal.NewConnection(internal.ConnectionID{}, conn, v, internal.ReuseConnection(v.config.IsConnectionReuse())), nil
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
		if conn != nil && v.tlsConfig != nil {
			conn = tls.Server(conn, v.tlsConfig)
		}
		if conn != nil && v.authConfig != nil {
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
			connErr.conn.Close()
		}
	}
	return nil
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_TCP, ListenTCP))
}
