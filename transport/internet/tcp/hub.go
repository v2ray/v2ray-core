package tcp

import (
	"context"
	gotls "crypto/tls"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

type TCPListener struct {
	listener   *net.TCPListener
	tlsConfig  *gotls.Config
	authConfig internet.ConnectionAuthenticator
	config     *Config
	addConn    internet.AddConnection
}

func ListenTCP(ctx context.Context, address net.Address, port net.Port, addConn internet.AddConnection) (internet.Listener, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	newError("listening TCP on ", address, ":", port).WriteToLog()
	networkSettings := internet.TransportSettingsFromContext(ctx)
	tcpSettings := networkSettings.(*Config)

	l := &TCPListener{
		listener: listener,
		config:   tcpSettings,
		addConn:  addConn,
	}

	if config := tls.ConfigFromContext(ctx, tls.WithNextProto("h2")); config != nil {
		l.tlsConfig = config.GetTLSConfig()
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
	go l.KeepAccepting(ctx)
	return l, nil
}

func (v *TCPListener) KeepAccepting(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
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
			newError("failed to accepted raw connections").Base(err).AtWarning().WriteToLog()
			continue
		}

		if v.tlsConfig != nil {
			conn = tls.Server(conn, v.tlsConfig)
		}
		if v.authConfig != nil {
			conn = v.authConfig.Server(conn)
		}

		v.addConn(context.Background(), internet.Connection(conn))
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
