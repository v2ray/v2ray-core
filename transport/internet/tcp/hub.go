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

type TCPListener struct {
	sync.Mutex
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
	log.Info("TCP|Listener: Listening on ", address, ":", port)
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

func (v *TCPListener) KeepAccepting() {
	for {
		select {
		case <-v.ctx.Done():
			return
		default:
		}
		conn, err := v.listener.Accept()
		v.Lock()
		if err != nil {
			log.Warning("TCP|Listener: Failed to accepted raw connections: ", err)
			v.Unlock()
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

		v.Unlock()
	}
}

func (v *TCPListener) Put(id internal.ConnectionID, conn net.Conn) {
	select {
	case <-v.ctx.Done():
		conn.Close()
		return
	case v.conns <- internal.NewConnection(internal.ConnectionID{}, conn, v, internal.ReuseConnection(v.config.IsConnectionReuse())):
	case <-time.After(time.Second * 5):
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
