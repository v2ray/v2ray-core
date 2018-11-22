package quic

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/transport/internet/tls"

	quic "github.com/lucas-clemente/quic-go"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

type clientSessions struct {
	access   sync.Mutex
	sessions map[net.Destination]quic.Session
}

func (s *clientSessions) getSession(destAddr net.Addr, config *Config, tlsConfig *tls.Config, sockopt *internet.SocketConfig) (quic.Session, error) {
	s.access.Lock()
	defer s.access.Unlock()

	if s.sessions == nil {
		s.sessions = make(map[net.Destination]quic.Session)
	}

	dest := net.DestinationFromAddr(destAddr)

	if session, found := s.sessions[dest]; found {
		select {
		case <-session.Context().Done():
			// Session has been closed. Creating a new one.
		default:
			return session, nil
		}
	}

	rawConn, err := internet.ListenSystemPacket(context.Background(), &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 0,
	}, sockopt)
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
		MaxIncomingUniStreams:                 -1,
	}

	conn, err := wrapSysConn(rawConn, config)
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	session, err := quic.DialContext(context.Background(), conn, destAddr, "", tlsConfig.GetTLSConfig(tls.WithDestination(dest)), quicConfig)
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	s.sessions[dest] = session
	return session, nil
}

var client clientSessions

func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	tlsConfig := tls.ConfigFromStreamSettings(streamSettings)
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			ServerName:    internalDomain,
			AllowInsecure: true,
		}
	}

	destAddr, err := net.ResolveUDPAddr("udp", dest.NetAddr())
	if err != nil {
		return nil, err
	}

	config := streamSettings.ProtocolSettings.(*Config)

	session, err := client.getSession(destAddr, config, tlsConfig, streamSettings.SocketSettings)
	if err != nil {
		return nil, err
	}

	conn, err := session.OpenStreamSync()
	if err != nil {
		return nil, err
	}

	return &interConn{
		stream: conn,
		local:  session.LocalAddr(),
		remote: destAddr,
	}, nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
