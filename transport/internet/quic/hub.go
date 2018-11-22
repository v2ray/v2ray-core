package quic

import (
	"context"
	"time"

	quic "github.com/lucas-clemente/quic-go"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/tls/cert"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

// Listener is an internet.Listener that listens for TCP connections.
type Listener struct {
	listener quic.Listener
	addConn  internet.ConnHandler
}

func (l *Listener) acceptStreams(session quic.Session) {
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			newError("failed to accept stream").Base(err).WriteToLog()
			session.Close()
			return
		}

		conn := &interConn{
			stream: stream,
			local:  session.LocalAddr(),
			remote: session.RemoteAddr(),
		}

		l.addConn(conn)
	}

}

func (l *Listener) keepAccepting() {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			newError("failed to accept QUIC sessions").Base(err).WriteToLog()
			l.listener.Close()
			return
		}
		go l.acceptStreams(conn)
	}
}

// Addr implements internet.Listener.Addr.
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

// Close implements internet.Listener.Close.
func (l *Listener) Close() error {
	return l.listener.Close()
}

// Listen creates a new Listener based on configurations.
func Listen(ctx context.Context, address net.Address, port net.Port, streamSettings *internet.MemoryStreamConfig, handler internet.ConnHandler) (internet.Listener, error) {
	if address.Family().IsDomain() {
		return nil, newError("domain address is not allows for listening quic")
	}

	tlsConfig := tls.ConfigFromStreamSettings(streamSettings)
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			Certificate: []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil, cert.DNSNames(internalDomain), cert.CommonName(internalDomain)))},
		}
	}

	config := streamSettings.ProtocolSettings.(*Config)
	rawConn, err := internet.ListenSystemPacket(context.Background(), &net.UDPAddr{
		IP:   address.IP(),
		Port: int(port),
	}, streamSettings.SocketSettings)

	if err != nil {
		return nil, err
	}

	quicConfig := &quic.Config{
		Versions:                              []quic.VersionNumber{quic.VersionMilestone0_10_0},
		ConnectionIDLength:                    12,
		KeepAlive:                             true,
		HandshakeTimeout:                      time.Second * 4,
		IdleTimeout:                           time.Second * 300,
		MaxReceiveStreamFlowControlWindow:     128 * 1024,
		MaxReceiveConnectionFlowControlWindow: 512 * 1024,
		MaxIncomingStreams:                    256,
		MaxIncomingUniStreams:                 -1,
	}

	conn, err := wrapSysConn(rawConn, config)
	if err != nil {
		conn.Close()
		return nil, err
	}

	qListener, err := quic.Listen(conn, tlsConfig.GetTLSConfig(), quicConfig)
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	listener := &Listener{
		listener: qListener,
		addConn:  handler,
	}

	go listener.keepAccepting()

	return listener, nil
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, Listen))
}
