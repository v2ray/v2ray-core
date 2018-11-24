package quic

import (
	"context"
	"sync"
	"time"

	quic "github.com/lucas-clemente/quic-go"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

type sessionContext struct {
	rawConn *sysConn
	session quic.Session
}

type clientSessions struct {
	access   sync.Mutex
	sessions map[net.Destination][]*sessionContext
}

func isActive(s quic.Session) bool {
	select {
	case <-s.Context().Done():
		return false
	default:
		return true
	}
}

func removeInactiveSessions(sessions []*sessionContext) []*sessionContext {
	activeSessions := make([]*sessionContext, 0, len(sessions))
	for _, s := range sessions {
		if isActive(s.session) {
			activeSessions = append(activeSessions, s)
		} else {
			s.rawConn.Close()
		}
	}

	if len(activeSessions) < len(sessions) {
		return activeSessions
	}

	return sessions
}

func openStream(sessions []*sessionContext) (quic.Stream, net.Addr) {
	for _, s := range sessions {
		if !isActive(s.session) {
			continue
		}

		stream, err := s.session.OpenStream()
		if err != nil {
			newError("failed to create stream").Base(err).AtWarning().WriteToLog()
			continue
		}
		return stream, s.session.LocalAddr()
	}

	return nil, nil
}

func (s *clientSessions) openConnection(destAddr net.Addr, config *Config, tlsConfig *tls.Config, sockopt *internet.SocketConfig) (internet.Connection, error) {
	s.access.Lock()
	defer s.access.Unlock()

	if s.sessions == nil {
		s.sessions = make(map[net.Destination][]*sessionContext)
	}

	dest := net.DestinationFromAddr(destAddr)

	var sessions []*sessionContext
	if s, found := s.sessions[dest]; found {
		sessions = s
	}

	{
		stream, local := openStream(sessions)
		if stream != nil {
			return &interConn{
				stream: stream,
				local:  local,
				remote: destAddr,
			}, nil
		}
	}

	sessions = removeInactiveSessions(sessions)

	rawConn, err := internet.ListenSystemPacket(context.Background(), &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 0,
	}, sockopt)
	if err != nil {
		return nil, err
	}

	quicConfig := &quic.Config{
		ConnectionIDLength:                    8,
		HandshakeTimeout:                      time.Second * 8,
		IdleTimeout:                           time.Second * 30,
		MaxReceiveStreamFlowControlWindow:     128 * 1024,
		MaxReceiveConnectionFlowControlWindow: 2 * 1024 * 1024,
		MaxIncomingUniStreams:                 -1,
		MaxIncomingStreams:                    32,
	}

	conn, err := wrapSysConn(rawConn, config)
	if err != nil {
		rawConn.Close()
		return nil, err
	}

	session, err := quic.DialContext(context.Background(), conn, destAddr, "", tlsConfig.GetTLSConfig(tls.WithDestination(dest)), quicConfig)
	if err != nil {
		conn.Close()
		return nil, err
	}

	s.sessions[dest] = append(sessions, &sessionContext{
		session: session,
		rawConn: conn,
	})
	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}
	return &interConn{
		stream: stream,
		local:  session.LocalAddr(),
		remote: destAddr,
	}, nil
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

	return client.openConnection(destAddr, config, tlsConfig, streamSettings.SocketSettings)
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
