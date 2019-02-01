package quic

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/task"
	quic "v2ray.com/core/external/github.com/lucas-clemente/quic-go"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

type sessionContext struct {
	rawConn *sysConn
	session quic.Session
}

var errSessionClosed = newError("session closed")

func (c *sessionContext) openStream(destAddr net.Addr) (*interConn, error) {
	if !isActive(c.session) {
		return nil, errSessionClosed
	}

	stream, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}

	conn := &interConn{
		stream: stream,
		local:  c.session.LocalAddr(),
		remote: destAddr,
	}

	return conn, nil
}

type clientSessions struct {
	access   sync.Mutex
	sessions map[net.Destination][]*sessionContext
	cleanup  *task.Periodic
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
			continue
		}
		if err := s.session.Close(); err != nil {
			newError("failed to close session").Base(err).WriteToLog()
		}
		if err := s.rawConn.Close(); err != nil {
			newError("failed to close raw connection").Base(err).WriteToLog()
		}
	}

	if len(activeSessions) < len(sessions) {
		return activeSessions
	}

	return sessions
}

func openStream(sessions []*sessionContext, destAddr net.Addr) *interConn {
	for _, s := range sessions {
		if !isActive(s.session) {
			continue
		}

		conn, err := s.openStream(destAddr)
		if err != nil {
			continue
		}

		return conn
	}

	return nil
}

func (s *clientSessions) cleanSessions() error {
	s.access.Lock()
	defer s.access.Unlock()

	if len(s.sessions) == 0 {
		return nil
	}

	newSessionMap := make(map[net.Destination][]*sessionContext)

	for dest, sessions := range s.sessions {
		sessions = removeInactiveSessions(sessions)
		if len(sessions) > 0 {
			newSessionMap[dest] = sessions
		}
	}

	s.sessions = newSessionMap
	return nil
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

	if true {
		conn := openStream(sessions, destAddr)
		if conn != nil {
			return conn, nil
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
		ConnectionIDLength: 12,
		HandshakeTimeout:   time.Second * 8,
		IdleTimeout:        time.Second * 30,
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

	context := &sessionContext{
		session: session,
		rawConn: conn,
	}
	s.sessions[dest] = append(sessions, context)
	return context.openStream(destAddr)
}

var client clientSessions

func init() {
	client.sessions = make(map[net.Destination][]*sessionContext)
	client.cleanup = &task.Periodic{
		Interval: time.Minute,
		Execute:  client.cleanSessions,
	}
	common.Must(client.cleanup.Start())
}

func Dial(ctx context.Context, dest net.Destination, streamSettings *internet.MemoryStreamConfig) (internet.Connection, error) {
	tlsConfig := tls.ConfigFromStreamSettings(streamSettings)
	if tlsConfig == nil {
		tlsConfig = &tls.Config{
			ServerName:    internalDomain,
			AllowInsecure: true,
		}
	}

	var destAddr *net.UDPAddr
	if dest.Address.Family().IsIP() {
		destAddr = &net.UDPAddr{
			IP:   dest.Address.IP(),
			Port: int(dest.Port),
		}
	} else {
		addr, err := net.ResolveUDPAddr("udp", dest.NetAddr())
		if err != nil {
			return nil, err
		}
		destAddr = addr
	}

	config := streamSettings.ProtocolSettings.(*Config)

	return client.openConnection(destAddr, config, tlsConfig, streamSettings.SocketSettings)
}
