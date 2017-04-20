package tcp

import (
	"context"
	gotls "crypto/tls"
	"net"
	"time"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

type TCPListener struct {
	ctx        context.Context
	listener   *net.TCPListener
	tlsConfig  *gotls.Config
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
	log.Trace(newError("listening TCP on ", address, ":", port))
	networkSettings := internet.TransportSettingsFromContext(ctx)
	tcpSettings := networkSettings.(*Config)

	l := &TCPListener{
		ctx:      ctx,
		listener: listener,
		config:   tcpSettings,
		conns:    conns,
	}
	if securitySettings := internet.SecuritySettingsFromContext(ctx); securitySettings != nil {
		tlsConfig, ok := securitySettings.(*tls.Config)
		if ok {
			l.tlsConfig = tlsConfig.GetTLSConfig()
		}
	}
	if tcpSettings.HeaderSettings != nil {
		headerConfig, err := tcpSettings.HeaderSettings.GetInstance()
		if err != nil {
			return nil, newError("invalid header settings").Base(err).AtError()
		}
		auth, err := internet.CreateConnectionAuthenticator(headerConfig)
		if err != nil {
			return nil, newError("invalid header settings.").Base(err).AtError()
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
			log.Trace(newError("failed to accepted raw connections").Base(err).AtWarning())
			continue
		}

		if v.tlsConfig != nil {
			conn = tls.Server(conn, v.tlsConfig)
		}
		if v.authConfig != nil {
			conn = v.authConfig.Server(conn)
		}

		select {
		case v.conns <- internet.Connection(conn):
		case <-time.After(time.Second * 5):
			conn.Close()
		}
	}
}

func (v *TCPListener) Addr() net.Addr {
	return v.listener.Addr()
}

func (v *TCPListener) Close() error {
	return v.listener.Close()
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_TCP, ListenTCP))
}
