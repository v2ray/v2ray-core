package tcp

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	ErrClosedListener = errors.New("Listener is closed.")
)

type TCPListener struct {
	ctx        context.Context
	listener   *net.TCPListener
	tlsConfig  *tls.Config
	authConfig internet.ConnectionAuthenticator
	config     *Config
	conns      chan<- internet.Connection
}

func ListenTCP(ctx context.Context, address v2net.Address, port v2net.Port, conns chan<- internet.Connection) (internet.Listener, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	log.Trace(errors.New("TCP|Listener: Listening on ", address, ":", port))
	networkSettings := internet.TransportSettingsFromContext(ctx)
	tcpSettings := networkSettings.(*Config)

	l := &TCPListener{
		ctx:      ctx,
		listener: listener,
		config:   tcpSettings,
		conns:    conns,
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
			return nil, errors.New("Internet|TCP: Invalid header settings.").Base(err)
		}
		auth, err := internet.CreateConnectionAuthenticator(headerConfig)
		if err != nil {
			return nil, errors.New("Internet|TCP: Invalid header settings.").Base(err)
		}
		l.authConfig = auth
	}
	go l.KeepAccepting()
	return l, nil
}

func (v *TCPListener) KeepAccepting() {
	for {
		select {
		case <-v.ctx.Done():
			return
		default:
		}
		var conn net.Conn
		err := retry.ExponentialBackoff(5, 200).On(func() error {
			rawConn, err := v.listener.Accept()
			if err != nil {
				return err
			}
			conn = rawConn
			return nil
		})
		if err != nil {
			log.Trace(errors.New("TCP|Listener: Failed to accepted raw connections: ", err).AtWarning())
			continue
		}

		if v.tlsConfig != nil {
			conn = tls.Server(conn, v.tlsConfig)
		}
		if v.authConfig != nil {
			conn = v.authConfig.Server(conn)
		}

		select {
		case v.conns <- internal.NewConnection(internal.ConnectionID{}, conn, v, internal.ReuseConnection(v.config.IsConnectionReuse())):
		case <-time.After(time.Second * 5):
			conn.Close()
		}
	}
}

func (v *TCPListener) Put(id internal.ConnectionID, conn net.Conn) {
	select {
	case v.conns <- internal.NewConnection(internal.ConnectionID{}, conn, v, internal.ReuseConnection(v.config.IsConnectionReuse())):
	case <-time.After(time.Second * 5):
		conn.Close()
	case <-v.ctx.Done():
		conn.Close()
	}
}

func (v *TCPListener) Addr() net.Addr {
	return v.listener.Addr()
}

func (v *TCPListener) Close() error {
	v.listener.Close()
	return nil
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_TCP, ListenTCP))
}
