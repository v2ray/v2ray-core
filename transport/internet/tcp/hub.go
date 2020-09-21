// +build !confonly

package tcp

import (
	"context"
	gotls "crypto/tls"
	"strings"
	"time"

	"github.com/pires/go-proxyproto"
	goxtls "github.com/xtls/go"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/xtls"
)

// Listener is an internet.Listener that listens for TCP connections.
type Listener struct {
	listener   net.Listener
	tlsConfig  *gotls.Config
	xtlsConfig *goxtls.Config
	authConfig internet.ConnectionAuthenticator
	config     *Config
	addConn    internet.ConnHandler
}

// ListenTCP creates a new Listener based on configurations.
func ListenTCP(ctx context.Context, address net.Address, port net.Port, streamSettings *internet.MemoryStreamConfig, handler internet.ConnHandler) (internet.Listener, error) {
	listener, err := internet.ListenSystem(ctx, &net.TCPAddr{
		IP:   address.IP(),
		Port: int(port),
	}, streamSettings.SocketSettings)
	if err != nil {
		return nil, newError("failed to listen TCP on", address, ":", port).Base(err)
	}
	newError("listening TCP on ", address, ":", port).WriteToLog(session.ExportIDToError(ctx))

	tcpSettings := streamSettings.ProtocolSettings.(*Config)
	var l *Listener

	if tcpSettings.AcceptProxyProtocol {
		policyFunc := func(upstream net.Addr) (proxyproto.Policy, error) { return proxyproto.REQUIRE, nil }
		l = &Listener{
			listener: &proxyproto.Listener{Listener: listener, Policy: policyFunc},
			config:   tcpSettings,
			addConn:  handler,
		}
		newError("accepting PROXY protocol").AtWarning().WriteToLog(session.ExportIDToError(ctx))
	} else {
		l = &Listener{
			listener: listener,
			config:   tcpSettings,
			addConn:  handler,
		}
	}

	if config := tls.ConfigFromStreamSettings(streamSettings); config != nil {
		l.tlsConfig = config.GetTLSConfig(tls.WithNextProto("h2"))
	}
	if config := xtls.ConfigFromStreamSettings(streamSettings); config != nil {
		l.xtlsConfig = config.GetXTLSConfig(xtls.WithNextProto("h2"))
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

	go l.keepAccepting()
	return l, nil
}

func (v *Listener) keepAccepting() {
	for {
		conn, err := v.listener.Accept()
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "closed") {
				break
			}
			newError("failed to accepted raw connections").Base(err).AtWarning().WriteToLog()
			if strings.Contains(errStr, "too many") {
				time.Sleep(time.Millisecond * 500)
			}
			continue
		}

		if v.tlsConfig != nil {
			conn = tls.Server(conn, v.tlsConfig)
		} else if v.xtlsConfig != nil {
			conn = xtls.Server(conn, v.xtlsConfig)
		}
		if v.authConfig != nil {
			conn = v.authConfig.Server(conn)
		}

		v.addConn(internet.Connection(conn))
	}
}

// Addr implements internet.Listener.Addr.
func (v *Listener) Addr() net.Addr {
	return v.listener.Addr()
}

// Close implements internet.Listener.Close.
func (v *Listener) Close() error {
	return v.listener.Close()
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, ListenTCP))
}
