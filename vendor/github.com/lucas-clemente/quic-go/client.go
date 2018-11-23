package quic

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

type client struct {
	mutex sync.Mutex

	conn connection
	// If the client is created with DialAddr, we create a packet conn.
	// If it is started with Dial, we take a packet conn as a parameter.
	createdPacketConn bool

	packetHandlers packetHandlerManager

	token []byte

	versionNegotiated                bool // has the server accepted our version
	receivedVersionNegotiationPacket bool
	negotiatedVersions               []protocol.VersionNumber // the list of versions from the version negotiation packet

	tlsConf *tls.Config
	config  *Config

	srcConnID      protocol.ConnectionID
	destConnID     protocol.ConnectionID
	origDestConnID protocol.ConnectionID // the destination conn ID used on the first Initial (before a Retry)

	initialVersion protocol.VersionNumber
	version        protocol.VersionNumber

	handshakeChan chan struct{}

	session quicSession

	logger utils.Logger
}

var _ packetHandler = &client{}

var (
	// make it possible to mock connection ID generation in the tests
	generateConnectionID           = protocol.GenerateConnectionID
	generateConnectionIDForInitial = protocol.GenerateConnectionIDForInitial
	errCloseSessionForNewVersion   = errors.New("closing session in order to recreate it with a new version")
	errCloseSessionForRetry        = errors.New("closing session in response to a stateless retry")
)

// DialAddr establishes a new QUIC connection to a server.
// The hostname for SNI is taken from the given address.
func DialAddr(
	addr string,
	tlsConf *tls.Config,
	config *Config,
) (Session, error) {
	return DialAddrContext(context.Background(), addr, tlsConf, config)
}

// DialAddrContext establishes a new QUIC connection to a server using the provided context.
// The hostname for SNI is taken from the given address.
func DialAddrContext(
	ctx context.Context,
	addr string,
	tlsConf *tls.Config,
	config *Config,
) (Session, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}
	return dialContext(ctx, udpConn, udpAddr, addr, tlsConf, config, true)
}

// Dial establishes a new QUIC connection to a server using a net.PacketConn.
// The host parameter is used for SNI.
func Dial(
	pconn net.PacketConn,
	remoteAddr net.Addr,
	host string,
	tlsConf *tls.Config,
	config *Config,
) (Session, error) {
	return DialContext(context.Background(), pconn, remoteAddr, host, tlsConf, config)
}

// DialContext establishes a new QUIC connection to a server using a net.PacketConn using the provided context.
// The host parameter is used for SNI.
func DialContext(
	ctx context.Context,
	pconn net.PacketConn,
	remoteAddr net.Addr,
	host string,
	tlsConf *tls.Config,
	config *Config,
) (Session, error) {
	return dialContext(ctx, pconn, remoteAddr, host, tlsConf, config, false)
}

func dialContext(
	ctx context.Context,
	pconn net.PacketConn,
	remoteAddr net.Addr,
	host string,
	tlsConf *tls.Config,
	config *Config,
	createdPacketConn bool,
) (Session, error) {
	config = populateClientConfig(config, createdPacketConn)
	packetHandlers, err := getMultiplexer().AddConn(pconn, config.ConnectionIDLength)
	if err != nil {
		return nil, err
	}
	c, err := newClient(pconn, remoteAddr, config, tlsConf, host, createdPacketConn)
	if err != nil {
		return nil, err
	}
	c.packetHandlers = packetHandlers
	if err := c.dial(ctx); err != nil {
		return nil, err
	}
	return c.session, nil
}

func newClient(
	pconn net.PacketConn,
	remoteAddr net.Addr,
	config *Config,
	tlsConf *tls.Config,
	host string,
	createdPacketConn bool,
) (*client, error) {
	if tlsConf == nil {
		tlsConf = &tls.Config{}
	}
	if tlsConf.ServerName == "" {
		var err error
		tlsConf.ServerName, _, err = net.SplitHostPort(host)
		if err != nil {
			return nil, err
		}
	}

	// check that all versions are actually supported
	if config != nil {
		for _, v := range config.Versions {
			if !protocol.IsValidVersion(v) {
				return nil, fmt.Errorf("%s is not a valid QUIC version", v)
			}
		}
	}
	c := &client{
		conn:              &conn{pconn: pconn, currentAddr: remoteAddr},
		createdPacketConn: createdPacketConn,
		tlsConf:           tlsConf,
		config:            config,
		version:           config.Versions[0],
		handshakeChan:     make(chan struct{}),
		logger:            utils.DefaultLogger.WithPrefix("client"),
	}
	return c, c.generateConnectionIDs()
}

// populateClientConfig populates fields in the quic.Config with their default values, if none are set
// it may be called with nil
func populateClientConfig(config *Config, createdPacketConn bool) *Config {
	if config == nil {
		config = &Config{}
	}
	versions := config.Versions
	if len(versions) == 0 {
		versions = protocol.SupportedVersions
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
	if connIDLen == 0 && !createdPacketConn {
		connIDLen = protocol.DefaultConnectionIDLength
	}

	return &Config{
		Versions:                              versions,
		HandshakeTimeout:                      handshakeTimeout,
		IdleTimeout:                           idleTimeout,
		ConnectionIDLength:                    connIDLen,
		MaxReceiveStreamFlowControlWindow:     maxReceiveStreamFlowControlWindow,
		MaxReceiveConnectionFlowControlWindow: maxReceiveConnectionFlowControlWindow,
		MaxIncomingStreams:                    maxIncomingStreams,
		MaxIncomingUniStreams:                 maxIncomingUniStreams,
		KeepAlive:                             config.KeepAlive,
	}
}

func (c *client) generateConnectionIDs() error {
	srcConnID, err := generateConnectionID(c.config.ConnectionIDLength)
	if err != nil {
		return err
	}
	destConnID, err := generateConnectionIDForInitial()
	if err != nil {
		return err
	}
	c.srcConnID = srcConnID
	c.destConnID = destConnID
	return nil
}

func (c *client) dial(ctx context.Context) error {
	c.logger.Infof("Starting new connection to %s (%s -> %s), source connection ID %s, destination connection ID %s, version %s", c.tlsConf.ServerName, c.conn.LocalAddr(), c.conn.RemoteAddr(), c.srcConnID, c.destConnID, c.version)

	if err := c.createNewTLSSession(c.version); err != nil {
		return err
	}
	err := c.establishSecureConnection(ctx)
	if err == errCloseSessionForRetry || err == errCloseSessionForNewVersion {
		return c.dial(ctx)
	}
	return err
}

// establishSecureConnection runs the session, and tries to establish a secure connection
// It returns:
// - errCloseSessionForNewVersion when the server sends a version negotiation packet
// - handshake.ErrCloseSessionForRetry when the server performs a stateless retry
// - any other error that might occur
// - when the connection is forward-secure
func (c *client) establishSecureConnection(ctx context.Context) error {
	errorChan := make(chan error, 1)

	go func() {
		err := c.session.run() // returns as soon as the session is closed
		if err != errCloseSessionForRetry && err != errCloseSessionForNewVersion && c.createdPacketConn {
			c.conn.Close()
		}
		errorChan <- err
	}()

	select {
	case <-ctx.Done():
		// The session will send a PeerGoingAway error to the server.
		c.session.Close()
		return ctx.Err()
	case err := <-errorChan:
		return err
	case <-c.handshakeChan:
		// handshake successfully completed
		return nil
	}
}

func (c *client) handlePacket(p *receivedPacket) {
	if err := c.handlePacketImpl(p); err != nil {
		c.logger.Errorf("error handling packet: %s", err)
	}
}

func (c *client) handlePacketImpl(p *receivedPacket) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// handle Version Negotiation Packets
	if p.header.IsVersionNegotiation {
		err := c.handleVersionNegotiationPacket(p.header)
		if err != nil {
			c.session.destroy(err)
		}
		// version negotiation packets have no payload
		return err
	}

	// reject packets with the wrong connection ID
	if !p.header.DestConnectionID.Equal(c.srcConnID) {
		return fmt.Errorf("received a packet with an unexpected connection ID (%s, expected %s)", p.header.DestConnectionID, c.srcConnID)
	}

	if p.header.Type == protocol.PacketTypeRetry {
		c.handleRetryPacket(p.header)
		return nil
	}

	// this is the first packet we are receiving
	// since it is not a Version Negotiation Packet, this means the server supports the suggested version
	if !c.versionNegotiated {
		c.versionNegotiated = true
	}

	c.session.handlePacket(p)
	return nil
}

func (c *client) handleVersionNegotiationPacket(hdr *wire.Header) error {
	// ignore delayed / duplicated version negotiation packets
	if c.receivedVersionNegotiationPacket || c.versionNegotiated {
		c.logger.Debugf("Received a delayed Version Negotiation Packet.")
		return nil
	}

	for _, v := range hdr.SupportedVersions {
		if v == c.version {
			// the version negotiation packet contains the version that we offered
			// this might be a packet sent by an attacker (or by a terribly broken server implementation)
			// ignore it
			return nil
		}
	}

	c.logger.Infof("Received a Version Negotiation Packet. Supported Versions: %s", hdr.SupportedVersions)
	newVersion, ok := protocol.ChooseSupportedVersion(c.config.Versions, hdr.SupportedVersions)
	if !ok {
		return qerr.InvalidVersion
	}
	c.receivedVersionNegotiationPacket = true
	c.negotiatedVersions = hdr.SupportedVersions

	// switch to negotiated version
	c.initialVersion = c.version
	c.version = newVersion
	if err := c.generateConnectionIDs(); err != nil {
		return err
	}

	c.logger.Infof("Switching to QUIC version %s. New connection ID: %s", newVersion, c.destConnID)
	c.session.destroy(errCloseSessionForNewVersion)
	return nil
}

func (c *client) handleRetryPacket(hdr *wire.Header) {
	c.logger.Debugf("<- Received Retry")
	hdr.Log(c.logger)
	if !hdr.OrigDestConnectionID.Equal(c.destConnID) {
		c.logger.Debugf("Ignoring spoofed Retry. Original Destination Connection ID: %s, expected: %s", hdr.OrigDestConnectionID, c.destConnID)
		return
	}
	if hdr.SrcConnectionID.Equal(c.destConnID) {
		c.logger.Debugf("Ignoring Retry, since the server didn't change the Source Connection ID.")
		return
	}
	// If a token is already set, this means that we already received a Retry from the server.
	// Ignore this Retry packet.
	if len(c.token) > 0 {
		c.logger.Debugf("Ignoring Retry, since a Retry was already received.")
		return
	}
	c.origDestConnID = c.destConnID
	c.destConnID = hdr.SrcConnectionID
	c.token = hdr.Token
	c.session.destroy(errCloseSessionForRetry)
}

func (c *client) createNewTLSSession(version protocol.VersionNumber) error {
	params := &handshake.TransportParameters{
		InitialMaxStreamDataBidiRemote: protocol.InitialMaxStreamData,
		InitialMaxStreamDataBidiLocal:  protocol.InitialMaxStreamData,
		InitialMaxStreamDataUni:        protocol.InitialMaxStreamData,
		InitialMaxData:                 protocol.InitialMaxData,
		IdleTimeout:                    c.config.IdleTimeout,
		MaxBidiStreams:                 uint64(c.config.MaxIncomingStreams),
		MaxUniStreams:                  uint64(c.config.MaxIncomingUniStreams),
		DisableMigration:               true,
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	runner := &runner{
		onHandshakeCompleteImpl: func(_ Session) { close(c.handshakeChan) },
		retireConnectionIDImpl:  c.packetHandlers.Retire,
		removeConnectionIDImpl:  c.packetHandlers.Remove,
	}
	sess, err := newClientSession(
		c.conn,
		runner,
		c.token,
		c.origDestConnID,
		c.destConnID,
		c.srcConnID,
		c.config,
		c.tlsConf,
		params,
		c.initialVersion,
		c.logger,
		c.version,
	)
	if err != nil {
		return err
	}
	c.session = sess
	c.packetHandlers.Add(c.srcConnID, c)
	return nil
}

func (c *client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.session == nil {
		return nil
	}
	return c.session.Close()
}

func (c *client) destroy(e error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.session == nil {
		return
	}
	c.session.destroy(e)
}

func (c *client) GetVersion() protocol.VersionNumber {
	c.mutex.Lock()
	v := c.version
	c.mutex.Unlock()
	return v
}

func (c *client) GetPerspective() protocol.Perspective {
	return protocol.PerspectiveClient
}
