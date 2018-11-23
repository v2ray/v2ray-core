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

type clientSessions struct {
	access   sync.Mutex
	sessions map[net.Destination][]quic.Session
}

func removeInactiveSessions(sessions []quic.Session) []quic.Session {
	lastActive := 0
	for _, s := range sessions {
		active := true
		select {
		case <-s.Context().Done():
			active = false
		default:
		}
		if active {
			sessions[lastActive] = s
			lastActive++
		}
	}

	if lastActive < len(sessions) {
		for i := 0; i < len(sessions); i++ {
			sessions[i] = nil
		}
		sessions = sessions[:lastActive]
	}

	return sessions
}

func openStream(sessions []quic.Session) (quic.Stream, net.Addr, error) {
	for _, s := range sessions {
		stream, err := s.OpenStream()
		if err != nil {
			newError("failed to create stream").Base(err).WriteToLog()
			continue
		}
		return stream, s.LocalAddr(), nil
	}

	return nil, nil, nil
}

func (s *clientSessions) openConnection(destAddr net.Addr, config *Config, tlsConfig *tls.Config, sockopt *internet.SocketConfig) (internet.Connection, error) {
	s.access.Lock()
	defer s.access.Unlock()

	if s.sessions == nil {
		s.sessions = make(map[net.Destination][]quic.Session)
	}

	dest := net.DestinationFromAddr(destAddr)

	var sessions []quic.Session
	if s, found := s.sessions[dest]; found {
		sessions = s
	}

	sessions = removeInactiveSessions(sessions)
	s.sessions[dest] = sessions

	stream, local, err := openStream(sessions)
	if err != nil {
		return nil, err
	}
	if stream != nil {
		return &interConn{
			stream: stream,
			local:  local,
			remote: destAddr,
		}, nil
	}

	rawConn, err := internet.ListenSystemPacket(context.Background(), &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 0,
	}, sockopt)
	if err != nil {
		return nil, err
	}

	quicConfig := &quic.Config{
		ConnectionIDLength:                    12,
		KeepAlive:                             true,
		HandshakeTimeout:                      time.Second * 4,
		IdleTimeout:                           time.Second * 60,
		MaxReceiveStreamFlowControlWindow:     256 * 1024,
		MaxReceiveConnectionFlowControlWindow: 2 * 1024 * 1024,
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

	s.sessions[dest] = append(sessions, session)
	stream, err = session.OpenStream()
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
