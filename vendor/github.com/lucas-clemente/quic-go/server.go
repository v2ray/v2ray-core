package quic

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

// packetHandler handles packets
type packetHandler interface {
	handlePacket(*receivedPacket)
	io.Closer
	destroy(error)
	GetPerspective() protocol.Perspective
}

type unknownPacketHandler interface {
	handlePacket(*receivedPacket)
	closeWithError(error) error
}

type packetHandlerManager interface {
	Add(protocol.ConnectionID, packetHandler)
	Retire(protocol.ConnectionID)
	Remove(protocol.ConnectionID)
	SetServer(unknownPacketHandler)
	CloseServer()
}

type quicSession interface {
	Session
	handlePacket(*receivedPacket)
	GetVersion() protocol.VersionNumber
	run() error
	destroy(error)
	closeRemote(error)
}

type sessionRunner interface {
	onHandshakeComplete(Session)
	retireConnectionID(protocol.ConnectionID)
	removeConnectionID(protocol.ConnectionID)
}

type runner struct {
	onHandshakeCompleteImpl func(Session)
	retireConnectionIDImpl  func(protocol.ConnectionID)
	removeConnectionIDImpl  func(protocol.ConnectionID)
}

func (r *runner) onHandshakeComplete(s Session)              { r.onHandshakeCompleteImpl(s) }
func (r *runner) retireConnectionID(c protocol.ConnectionID) { r.retireConnectionIDImpl(c) }
func (r *runner) removeConnectionID(c protocol.ConnectionID) { r.removeConnectionIDImpl(c) }

var _ sessionRunner = &runner{}

// A Listener of QUIC
type server struct {
	mutex sync.Mutex

	tlsConf *tls.Config
	config  *Config

	conn net.PacketConn
	// If the server is started with ListenAddr, we create a packet conn.
	// If it is started with Listen, we take a packet conn as a parameter.
	createdPacketConn bool

	cookieGenerator *handshake.CookieGenerator

	sessionHandler packetHandlerManager

	// set as a member, so they can be set in the tests
	newSession func(connection, sessionRunner, protocol.ConnectionID /* original connection ID */, protocol.ConnectionID /* destination connection ID */, protocol.ConnectionID /* source connection ID */, *Config, *tls.Config, *handshake.TransportParameters, utils.Logger, protocol.VersionNumber) (quicSession, error)

	serverError error
	errorChan   chan struct{}
	closed      bool

	sessionQueue chan Session

	sessionRunner sessionRunner

	logger utils.Logger
}

var _ Listener = &server{}
var _ unknownPacketHandler = &server{}

// ListenAddr creates a QUIC server listening on a given address.
// The tls.Config must not be nil and must contain a certificate configuration.
// The quic.Config may be nil, in that case the default values will be used.
func ListenAddr(addr string, tlsConf *tls.Config, config *Config) (Listener, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	serv, err := listen(conn, tlsConf, config)
	if err != nil {
		return nil, err
	}
	serv.createdPacketConn = true
	return serv, nil
}

// Listen listens for QUIC connections on a given net.PacketConn.
// A single PacketConn only be used for a single call to Listen.
// The PacketConn can be used for simultaneous calls to Dial.
// QUIC connection IDs are used for demultiplexing the different connections.
// The tls.Config must not be nil and must contain a certificate configuration.
// The quic.Config may be nil, in that case the default values will be used.
func Listen(conn net.PacketConn, tlsConf *tls.Config, config *Config) (Listener, error) {
	return listen(conn, tlsConf, config)
}

func listen(conn net.PacketConn, tlsConf *tls.Config, config *Config) (*server, error) {
	config = populateServerConfig(config)
	for _, v := range config.Versions {
		if !protocol.IsValidVersion(v) {
			return nil, fmt.Errorf("%s is not a valid QUIC version", v)
		}
	}

	sessionHandler, err := getMultiplexer().AddConn(conn, config.ConnectionIDLength)
	if err != nil {
		return nil, err
	}
	s := &server{
		conn:           conn,
		tlsConf:        tlsConf,
		config:         config,
		sessionHandler: sessionHandler,
		sessionQueue:   make(chan Session, 5),
		errorChan:      make(chan struct{}),
		newSession:     newSession,
		logger:         utils.DefaultLogger.WithPrefix("server"),
	}
	if err := s.setup(); err != nil {
		return nil, err
	}
	sessionHandler.SetServer(s)
	s.logger.Debugf("Listening for %s connections on %s", conn.LocalAddr().Network(), conn.LocalAddr().String())
	return s, nil
}

func (s *server) setup() error {
	s.sessionRunner = &runner{
		onHandshakeCompleteImpl: func(sess Session) { s.sessionQueue <- sess },
		retireConnectionIDImpl:  s.sessionHandler.Retire,
		removeConnectionIDImpl:  s.sessionHandler.Remove,
	}
	cookieGenerator, err := handshake.NewCookieGenerator()
	if err != nil {
		return err
	}
	s.cookieGenerator = cookieGenerator
	return nil
}

var defaultAcceptCookie = func(clientAddr net.Addr, cookie *Cookie) bool {
	if cookie == nil {
		return false
	}
	if time.Now().After(cookie.SentTime.Add(protocol.CookieExpiryTime)) {
		return false
	}
	var sourceAddr string
	if udpAddr, ok := clientAddr.(*net.UDPAddr); ok {
		sourceAddr = udpAddr.IP.String()
	} else {
		sourceAddr = clientAddr.String()
	}
	return sourceAddr == cookie.RemoteAddr
}

// populateServerConfig populates fields in the quic.Config with their default values, if none are set
// it may be called with nil
func populateServerConfig(config *Config) *Config {
	if config == nil {
		config = &Config{}
	}
	versions := config.Versions
	if len(versions) == 0 {
		versions = protocol.SupportedVersions
	}

	vsa := defaultAcceptCookie
	if config.AcceptCookie != nil {
		vsa = config.AcceptCookie
	}

	handshakeTimeout := protocol.DefaultHandshakeTimeout
	if config.HandshakeTimeout != 0 {
		handshakeTimeout = config.HandshakeTimeout
	}
	idleTimeout := protocol.DefaultIdleTimeout
	if config.IdleTimeout != 0 {
		idleTimeout = config.IdleTimeout
	}

	maxReceiveStreamFlowControlWindow := config.MaxReceiveStreamFlowControlWindow
	if maxReceiveStreamFlowControlWindow == 0 {
		maxReceiveStreamFlowControlWindow = protocol.DefaultMaxReceiveStreamFlowControlWindow
	}
	maxReceiveConnectionFlowControlWindow := config.MaxReceiveConnectionFlowControlWindow
	if maxReceiveConnectionFlowControlWindow == 0 {
		maxReceiveConnectionFlowControlWindow = protocol.DefaultMaxReceiveConnectionFlowControlWindow
	}
	maxIncomingStreams := config.MaxIncomingStreams
	if maxIncomingStreams == 0 {
		maxIncomingStreams = protocol.DefaultMaxIncomingStreams
	} else if maxIncomingStreams < 0 {
		maxIncomingStreams = 0
	}
	maxIncomingUniStreams := config.MaxIncomingUniStreams
	if maxIncomingUniStreams == 0 {
		maxIncomingUniStreams = protocol.DefaultMaxIncomingUniStreams
	} else if maxIncomingUniStreams < 0 {
		maxIncomingUniStreams = 0
	}
	connIDLen := config.ConnectionIDLength
	if connIDLen == 0 {
		connIDLen = protocol.DefaultConnectionIDLength
	}

	return &Config{
		Versions:                              versions,
		HandshakeTimeout:                      handshakeTimeout,
		IdleTimeout:                           idleTimeout,
		AcceptCookie:                          vsa,
		KeepAlive:                             config.KeepAlive,
		MaxReceiveStreamFlowControlWindow:     maxReceiveStreamFlowControlWindow,
		MaxReceiveConnectionFlowControlWindow: maxReceiveConnectionFlowControlWindow,
		MaxIncomingStreams:                    maxIncomingStreams,
		MaxIncomingUniStreams:                 maxIncomingUniStreams,
		ConnectionIDLength:                    connIDLen,
	}
}

// Accept returns newly openend sessions
func (s *server) Accept() (Session, error) {
	var sess Session
	select {
	case sess = <-s.sessionQueue:
		return sess, nil
	case <-s.errorChan:
		return nil, s.serverError
	}
}

// Close the server
func (s *server) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return nil
	}
	return s.closeWithMutex()
}

func (s *server) closeWithMutex() error {
	s.sessionHandler.CloseServer()
	if s.serverError == nil {
		s.serverError = errors.New("server closed")
	}
	var err error
	// If the server was started with ListenAddr, we created the packet conn.
	// We need to close it in order to make the go routine reading from that conn return.
	if s.createdPacketConn {
		err = s.conn.Close()
	}
	s.closed = true
	close(s.errorChan)
	return err
}

func (s *server) closeWithError(e error) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.closed {
		return nil
	}
	s.serverError = e
	return s.closeWithMutex()
}

// Addr returns the server's network address
func (s *server) Addr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *server) handlePacket(p *receivedPacket) {
	hdr := p.hdr

	// send a Version Negotiation Packet if the client is speaking a different protocol version
	if !protocol.IsSupportedVersion(s.config.Versions, hdr.Version) {
		go s.sendVersionNegotiationPacket(p)
		return
	}
	if hdr.Type == protocol.PacketTypeInitial {
		go s.handleInitial(p)
	}
	// TODO(#943): send Stateless Reset
}

func (s *server) handleInitial(p *receivedPacket) {
	// TODO: add a check that DestConnID == SrcConnID
	s.logger.Debugf("<- Received Initial packet.")
	sess, connID, err := s.handleInitialImpl(p)
	if err != nil {
		s.logger.Errorf("Error occurred handling initial packet: %s", err)
		return
	}
	if sess == nil { // a retry was done
		return
	}
	serverSession := newServerSession(sess, s.config, s.logger)
	s.sessionHandler.Add(connID, serverSession)
}

func (s *server) handleInitialImpl(p *receivedPacket) (quicSession, protocol.ConnectionID, error) {
	hdr := p.hdr
	if len(hdr.Token) == 0 && hdr.DestConnectionID.Len() < protocol.MinConnectionIDLenInitial {
		return nil, nil, errors.New("dropping Initial packet with too short connection ID")
	}
	if len(p.data) < protocol.MinInitialPacketSize {
		return nil, nil, errors.New("dropping too small Initial packet")
	}

	var cookie *Cookie
	var origDestConnectionID protocol.ConnectionID
	if len(hdr.Token) > 0 {
		c, err := s.cookieGenerator.DecodeToken(hdr.Token)
		if err == nil {
			cookie = &Cookie{
				RemoteAddr: c.RemoteAddr,
				SentTime:   c.SentTime,
			}
			origDestConnectionID = c.OriginalDestConnectionID
		}
	}
	if !s.config.AcceptCookie(p.remoteAddr, cookie) {
		// Log the Initial packet now.
		// If no Retry is sent, the packet will be logged by the session.
		(&wire.ExtendedHeader{Header: *p.hdr}).Log(s.logger)
		return nil, nil, s.sendRetry(p.remoteAddr, hdr)
	}

	connID, err := protocol.GenerateConnectionID(s.config.ConnectionIDLength)
	if err != nil {
		return nil, nil, err
	}
	s.logger.Debugf("Changing connection ID to %s.", connID)
	sess, err := s.createNewSession(
		p.remoteAddr,
		origDestConnectionID,
		hdr.DestConnectionID,
		hdr.SrcConnectionID,
		connID,
		hdr.Version,
	)
	if err != nil {
		return nil, nil, err
	}
	sess.handlePacket(p)
	return sess, connID, nil
}

func (s *server) createNewSession(
	remoteAddr net.Addr,
	origDestConnID protocol.ConnectionID,
	clientDestConnID protocol.ConnectionID,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	version protocol.VersionNumber,
) (quicSession, error) {
	params := &handshake.TransportParameters{
		InitialMaxStreamDataBidiLocal:  protocol.InitialMaxStreamData,
		InitialMaxStreamDataBidiRemote: protocol.InitialMaxStreamData,
		InitialMaxStreamDataUni:        protocol.InitialMaxStreamData,
		InitialMaxData:                 protocol.InitialMaxData,
		IdleTimeout:                    s.config.IdleTimeout,
		MaxBidiStreams:                 uint64(s.config.MaxIncomingStreams),
		MaxUniStreams:                  uint64(s.config.MaxIncomingUniStreams),
		DisableMigration:               true,
		// TODO(#855): generate a real token
		StatelessResetToken:  bytes.Repeat([]byte{42}, 16),
		OriginalConnectionID: origDestConnID,
	}
	sess, err := s.newSession(
		&conn{pconn: s.conn, currentAddr: remoteAddr},
		s.sessionRunner,
		clientDestConnID,
		destConnID,
		srcConnID,
		s.config,
		s.tlsConf,
		params,
		s.logger,
		version,
	)
	if err != nil {
		return nil, err
	}
	go sess.run()
	return sess, nil
}

func (s *server) sendRetry(remoteAddr net.Addr, hdr *wire.Header) error {
	token, err := s.cookieGenerator.NewToken(remoteAddr, hdr.DestConnectionID)
	if err != nil {
		return err
	}
	connID, err := protocol.GenerateConnectionID(s.config.ConnectionIDLength)
	if err != nil {
		return err
	}
	replyHdr := &wire.ExtendedHeader{}
	replyHdr.IsLongHeader = true
	replyHdr.Type = protocol.PacketTypeRetry
	replyHdr.Version = hdr.Version
	replyHdr.SrcConnectionID = connID
	replyHdr.DestConnectionID = hdr.SrcConnectionID
	replyHdr.OrigDestConnectionID = hdr.DestConnectionID
	replyHdr.Token = token
	s.logger.Debugf("Changing connection ID to %s.\n-> Sending Retry", connID)
	replyHdr.Log(s.logger)
	buf := &bytes.Buffer{}
	if err := replyHdr.Write(buf, hdr.Version); err != nil {
		return err
	}
	if _, err := s.conn.WriteTo(buf.Bytes(), remoteAddr); err != nil {
		s.logger.Debugf("Error sending Retry: %s", err)
	}
	return nil
}

func (s *server) sendVersionNegotiationPacket(p *receivedPacket) {
	hdr := p.hdr
	s.logger.Debugf("Client offered version %s, sending Version Negotiation", hdr.Version)
	data, err := wire.ComposeVersionNegotiation(hdr.SrcConnectionID, hdr.DestConnectionID, s.config.Versions)
	if err != nil {
		s.logger.Debugf("Error composing Version Negotiation: %s", err)
		return
	}
	if _, err := s.conn.WriteTo(data, p.remoteAddr); err != nil {
		s.logger.Debugf("Error sending Version Negotiation: %s", err)
	}
}
