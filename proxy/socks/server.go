package socks

import (
	"io"
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bufio"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/socks/protocol"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

var (
	ErrUnsupportedSocksCommand = errors.New("Unsupported socks command.")
	ErrUnsupportedAuthMethod   = errors.New("Unsupported auth method.")
)

// Server is a SOCKS 5 proxy server
type Server struct {
	tcpMutex         sync.RWMutex
	udpMutex         sync.RWMutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *ServerConfig
	tcpListener      *internet.TCPHub
	udpHub           *udp.UDPHub
	udpAddress       v2net.Destination
	udpServer        *udp.UDPServer
	meta             *proxy.InboundHandlerMeta
}

// NewServer creates a new Server object.
func NewServer(config *ServerConfig, space app.Space, meta *proxy.InboundHandlerMeta) *Server {
	s := &Server{
		config: config,
		meta:   meta,
	}
	space.InitializeApplication(func() error {
		if !space.HasApp(dispatcher.APP_ID) {
			return errors.New("Socks|Server: Dispatcher is not found in the space.")
		}
		s.packetDispatcher = space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
		return nil
	})
	return s
}

// Port implements InboundHandler.Port().
func (v *Server) Port() v2net.Port {
	return v.meta.Port
}

// Close implements InboundHandler.Close().
func (v *Server) Close() {
	v.accepting = false
	if v.tcpListener != nil {
		v.tcpMutex.Lock()
		v.tcpListener.Close()
		v.tcpListener = nil
		v.tcpMutex.Unlock()
	}
	if v.udpHub != nil {
		v.udpMutex.Lock()
		v.udpHub.Close()
		v.udpHub = nil
		v.udpMutex.Unlock()
	}
}

// Listen implements InboundHandler.Listen().
func (v *Server) Start() error {
	if v.accepting {
		return nil
	}

	listener, err := internet.ListenTCP(
		v.meta.Address,
		v.meta.Port,
		v.handleConnection,
		v.meta.StreamSettings)
	if err != nil {
		log.Error("Socks: failed to listen on ", v.meta.Address, ":", v.meta.Port, ": ", err)
		return err
	}
	v.accepting = true
	v.tcpMutex.Lock()
	v.tcpListener = listener
	v.tcpMutex.Unlock()
	if v.config.UdpEnabled {
		v.listenUDP()
	}
	return nil
}

func (v *Server) handleConnection(connection internet.Connection) {
	defer connection.Close()

	timedReader := v2net.NewTimeOutReader(v.config.Timeout, connection)
	reader := bufio.NewReader(timedReader)
	defer reader.Release()

	writer := bufio.NewWriter(connection)
	defer writer.Release()

	auth, auth4, err := protocol.ReadAuthentication(reader)
	if err != nil && errors.Cause(err) != protocol.Socks4Downgrade {
		if errors.Cause(err) != io.EOF {
			log.Warning("Socks: failed to read authentication: ", err)
		}
		return
	}

	clientAddr := v2net.DestinationFromAddr(connection.RemoteAddr())
	if err != nil && err == protocol.Socks4Downgrade {
		v.handleSocks4(clientAddr, reader, writer, auth4)
	} else {
		v.handleSocks5(clientAddr, reader, writer, auth)
	}
}

func (v *Server) handleSocks5(clientAddr v2net.Destination, reader *bufio.BufferedReader, writer *bufio.BufferedWriter, auth protocol.Socks5AuthenticationRequest) error {
	expectedAuthMethod := protocol.AuthNotRequired
	if v.config.AuthType == AuthType_PASSWORD {
		expectedAuthMethod = protocol.AuthUserPass
	}

	if !auth.HasAuthMethod(expectedAuthMethod) {
		authResponse := protocol.NewAuthenticationResponse(protocol.AuthNoMatchingMethod)
		err := protocol.WriteAuthentication(writer, authResponse)
		writer.Flush()
		if err != nil {
			log.Warning("Socks: failed to write authentication: ", err)
			return err
		}
		log.Warning("Socks: client doesn't support any allowed auth methods.")
		return ErrUnsupportedAuthMethod
	}

	authResponse := protocol.NewAuthenticationResponse(expectedAuthMethod)
	protocol.WriteAuthentication(writer, authResponse)
	err := writer.Flush()
	if err != nil {
		log.Error("Socks: failed to write authentication: ", err)
		return err
	}
	if v.config.AuthType == AuthType_PASSWORD {
		upRequest, err := protocol.ReadUserPassRequest(reader)
		if err != nil {
			log.Warning("Socks: failed to read username and password: ", err)
			return err
		}
		status := byte(0)
		if !v.config.HasAccount(upRequest.Username(), upRequest.Password()) {
			status = byte(0xFF)
		}
		upResponse := protocol.NewSocks5UserPassResponse(status)
		err = protocol.WriteUserPassResponse(writer, upResponse)
		writer.Flush()
		if err != nil {
			log.Error("Socks: failed to write user pass response: ", err)
			return err
		}
		if status != byte(0) {
			log.Warning("Socks: Invalid user account: ", upRequest.AuthDetail())
			log.Access(clientAddr, "", log.AccessRejected, crypto.ErrAuthenticationFailed)
			return crypto.ErrAuthenticationFailed
		}
	}

	request, err := protocol.ReadRequest(reader)
	if err != nil {
		log.Warning("Socks: failed to read request: ", err)
		return err
	}

	if request.Command == protocol.CmdUdpAssociate && v.config.UdpEnabled {
		return v.handleUDP(reader, writer)
	}

	if request.Command == protocol.CmdBind || request.Command == protocol.CmdUdpAssociate {
		response := protocol.NewSocks5Response()
		response.Error = protocol.ErrorCommandNotSupported
		response.Port = v2net.Port(0)
		response.SetIPv4([]byte{0, 0, 0, 0})

		response.Write(writer)
		writer.Flush()
		if err != nil {
			log.Error("Socks: failed to write response: ", err)
			return err
		}
		log.Warning("Socks: Unsupported socks command ", request.Command)
		return ErrUnsupportedSocksCommand
	}

	response := protocol.NewSocks5Response()
	response.Error = protocol.ErrorSuccess

	// Some SOCKS software requires a value other than dest. Let's fake one:
	response.Port = v2net.Port(1717)
	response.SetIPv4([]byte{0, 0, 0, 0})

	response.Write(writer)
	if err != nil {
		log.Error("Socks: failed to write response: ", err)
		return err
	}

	reader.SetCached(false)
	writer.SetCached(false)

	dest := request.Destination()
	session := &proxy.SessionInfo{
		Source:      clientAddr,
		Destination: dest,
		Inbound:     v.meta,
	}
	log.Info("Socks: TCP Connect request to ", dest)
	log.Access(clientAddr, dest, log.AccessAccepted, "")

	v.transport(reader, writer, session)
	return nil
}

func (v *Server) handleUDP(reader io.Reader, writer *bufio.BufferedWriter) error {
	response := protocol.NewSocks5Response()
	response.Error = protocol.ErrorSuccess

	udpAddr := v.udpAddress

	response.Port = udpAddr.Port
	switch udpAddr.Address.Family() {
	case v2net.AddressFamilyIPv4:
		response.SetIPv4(udpAddr.Address.IP())
	case v2net.AddressFamilyIPv6:
		response.SetIPv6(udpAddr.Address.IP())
	case v2net.AddressFamilyDomain:
		response.SetDomain(udpAddr.Address.Domain())
	}

	response.Write(writer)
	err := writer.Flush()

	if err != nil {
		log.Error("Socks: failed to write response: ", err)
		return err
	}

	// The TCP connection closes after v method returns. We need to wait until
	// the client closes it.
	// TODO: get notified from UDP part
	<-time.After(5 * time.Minute)

	return nil
}

func (v *Server) handleSocks4(clientAddr v2net.Destination, reader *bufio.BufferedReader, writer *bufio.BufferedWriter, auth protocol.Socks4AuthenticationRequest) error {
	result := protocol.Socks4RequestGranted
	if auth.Command == protocol.CmdBind {
		result = protocol.Socks4RequestRejected
	}
	socks4Response := protocol.NewSocks4AuthenticationResponse(result, auth.Port, auth.IP[:])

	socks4Response.Write(writer)

	if result == protocol.Socks4RequestRejected {
		log.Warning("Socks: Unsupported socks 4 command ", auth.Command)
		log.Access(clientAddr, "", log.AccessRejected, ErrUnsupportedSocksCommand)
		return ErrUnsupportedSocksCommand
	}

	reader.SetCached(false)
	writer.SetCached(false)

	dest := v2net.TCPDestination(v2net.IPAddress(auth.IP[:]), auth.Port)
	session := &proxy.SessionInfo{
		Source:      clientAddr,
		Destination: dest,
		Inbound:     v.meta,
	}
	log.Access(clientAddr, dest, log.AccessAccepted, "")
	v.transport(reader, writer, session)
	return nil
}

func (v *Server) transport(reader io.Reader, writer io.Writer, session *proxy.SessionInfo) {
	ray := v.packetDispatcher.DispatchToOutbound(session)
	input := ray.InboundInput()
	output := ray.InboundOutput()

	defer input.Close()
	defer output.Release()

	go func() {
		v2reader := buf.NewReader(reader)
		defer v2reader.Release()

		if err := buf.PipeUntilEOF(v2reader, input); err != nil {
			log.Info("Socks|Server: Failed to transport all TCP request: ", err)
		}
		input.Close()
	}()

	v2writer := buf.NewWriter(writer)
	defer v2writer.Release()

	if err := buf.PipeUntilEOF(output, v2writer); err != nil {
		log.Info("Socks|Server: Failed to transport all TCP response: ", err)
	}
	output.Release()
}

type ServerFactory struct{}

func (v *ServerFactory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_RawTCP},
	}
}

func (v *ServerFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	return NewServer(rawConfig.(*ServerConfig), space, meta), nil
}

func init() {
	proxy.MustRegisterInboundHandlerCreator(serial.GetMessageType(new(ServerConfig)), new(ServerFactory))
}
