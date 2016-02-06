package socks

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
	"github.com/v2ray/v2ray-core/transport/hub"
)

var (
	ErrorUnsupportedSocksCommand = errors.New("Unsupported socks command.")
	ErrorUnsupportedAuthMethod   = errors.New("Unsupported auth method.")
)

// SocksServer is a SOCKS 5 proxy server
type SocksServer struct {
	tcpMutex         sync.RWMutex
	udpMutex         sync.RWMutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *Config
	tcpListener      *hub.TCPHub
	udpHub           *hub.UDPHub
	udpAddress       v2net.Destination
	udpServer        *hub.UDPServer
	listeningPort    v2net.Port
}

// NewSocksSocks creates a new SocksServer object.
func NewSocksServer(config *Config, packetDispatcher dispatcher.PacketDispatcher) *SocksServer {
	return &SocksServer{
		config:           config,
		packetDispatcher: packetDispatcher,
	}
}

// Port implements InboundHandler.Port().
func (this *SocksServer) Port() v2net.Port {
	return this.listeningPort
}

// Close implements InboundHandler.Close().
func (this *SocksServer) Close() {
	this.accepting = false
	if this.tcpListener != nil {
		this.tcpMutex.Lock()
		this.tcpListener.Close()
		this.tcpListener = nil
		this.tcpMutex.Unlock()
	}
	if this.udpHub != nil {
		this.udpMutex.Lock()
		this.udpHub.Close()
		this.udpHub = nil
		this.udpMutex.Unlock()
	}
}

// Listen implements InboundHandler.Listen().
func (this *SocksServer) Listen(port v2net.Port) error {
	if this.accepting {
		if this.listeningPort == port {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}
	this.listeningPort = port

	listener, err := hub.ListenTCP(port, this.handleConnection)
	if err != nil {
		log.Error("Socks: failed to listen on port ", port, ": ", err)
		return err
	}
	this.accepting = true
	this.tcpMutex.Lock()
	this.tcpListener = listener
	this.tcpMutex.Unlock()
	if this.config.UDPEnabled {
		this.listenUDP(port)
	}
	return nil
}

func (this *SocksServer) handleConnection(connection *hub.TCPConn) {
	defer connection.Close()

	reader := v2net.NewTimeOutReader(120, connection)

	auth, auth4, err := protocol.ReadAuthentication(reader)
	if err != nil && err != protocol.Socks4Downgrade {
		log.Error("Socks: failed to read authentication: ", err)
		return
	}

	if err != nil && err == protocol.Socks4Downgrade {
		this.handleSocks4(reader, connection, auth4)
	} else {
		this.handleSocks5(reader, connection, auth)
	}
}

func (this *SocksServer) handleSocks5(reader *v2net.TimeOutReader, writer io.Writer, auth protocol.Socks5AuthenticationRequest) error {
	expectedAuthMethod := protocol.AuthNotRequired
	if this.config.AuthType == AuthTypePassword {
		expectedAuthMethod = protocol.AuthUserPass
	}

	if !auth.HasAuthMethod(expectedAuthMethod) {
		authResponse := protocol.NewAuthenticationResponse(protocol.AuthNoMatchingMethod)
		err := protocol.WriteAuthentication(writer, authResponse)
		if err != nil {
			log.Error("Socks: failed to write authentication: ", err)
			return err
		}
		log.Warning("Socks: client doesn't support any allowed auth methods.")
		return ErrorUnsupportedAuthMethod
	}

	authResponse := protocol.NewAuthenticationResponse(expectedAuthMethod)
	err := protocol.WriteAuthentication(writer, authResponse)
	if err != nil {
		log.Error("Socks: failed to write authentication: ", err)
		return err
	}
	if this.config.AuthType == AuthTypePassword {
		upRequest, err := protocol.ReadUserPassRequest(reader)
		if err != nil {
			log.Error("Socks: failed to read username and password: ", err)
			return err
		}
		status := byte(0)
		if !this.config.HasAccount(upRequest.Username(), upRequest.Password()) {
			status = byte(0xFF)
		}
		upResponse := protocol.NewSocks5UserPassResponse(status)
		err = protocol.WriteUserPassResponse(writer, upResponse)
		if err != nil {
			log.Error("Socks: failed to write user pass response: ", err)
			return err
		}
		if status != byte(0) {
			log.Warning("Socks: Invalid user account: ", upRequest.AuthDetail())
			return proxy.ErrorInvalidAuthentication
		}
	}

	request, err := protocol.ReadRequest(reader)
	if err != nil {
		log.Error("Socks: failed to read request: ", err)
		return err
	}

	if request.Command == protocol.CmdUdpAssociate && this.config.UDPEnabled {
		return this.handleUDP(reader, writer)
	}

	if request.Command == protocol.CmdBind || request.Command == protocol.CmdUdpAssociate {
		response := protocol.NewSocks5Response()
		response.Error = protocol.ErrorCommandNotSupported
		response.Port = v2net.Port(0)
		response.SetIPv4([]byte{0, 0, 0, 0})

		responseBuffer := alloc.NewSmallBuffer().Clear()
		response.Write(responseBuffer)
		_, err = writer.Write(responseBuffer.Value)
		responseBuffer.Release()
		if err != nil {
			log.Error("Socks: failed to write response: ", err)
			return err
		}
		log.Warning("Socks: Unsupported socks command ", request.Command)
		return ErrorUnsupportedSocksCommand
	}

	response := protocol.NewSocks5Response()
	response.Error = protocol.ErrorSuccess

	// Some SOCKS software requires a value other than dest. Let's fake one:
	response.Port = v2net.Port(1717)
	response.SetIPv4([]byte{0, 0, 0, 0})

	responseBuffer := alloc.NewSmallBuffer().Clear()
	response.Write(responseBuffer)
	_, err = writer.Write(responseBuffer.Value)
	responseBuffer.Release()
	if err != nil {
		log.Error("Socks: failed to write response: ", err)
		return err
	}

	dest := request.Destination()
	log.Info("Socks: TCP Connect request to ", dest)

	packet := v2net.NewPacket(dest, nil, true)
	this.transport(reader, writer, packet)
	return nil
}

func (this *SocksServer) handleUDP(reader *v2net.TimeOutReader, writer io.Writer) error {
	response := protocol.NewSocks5Response()
	response.Error = protocol.ErrorSuccess

	udpAddr := this.udpAddress

	response.Port = udpAddr.Port()
	switch {
	case udpAddr.Address().IsIPv4():
		response.SetIPv4(udpAddr.Address().IP())
	case udpAddr.Address().IsIPv6():
		response.SetIPv6(udpAddr.Address().IP())
	case udpAddr.Address().IsDomain():
		response.SetDomain(udpAddr.Address().Domain())
	}

	responseBuffer := alloc.NewSmallBuffer().Clear()
	response.Write(responseBuffer)
	_, err := writer.Write(responseBuffer.Value)
	responseBuffer.Release()

	if err != nil {
		log.Error("Socks: failed to write response: ", err)
		return err
	}

	reader.SetTimeOut(300)     /* 5 minutes */
	v2io.ReadFrom(reader, nil) // Just in case of anything left in the socket
	// The TCP connection closes after this method returns. We need to wait until
	// the client closes it.
	// TODO: get notified from UDP part
	<-time.After(5 * time.Minute)

	return nil
}

func (this *SocksServer) handleSocks4(reader io.Reader, writer io.Writer, auth protocol.Socks4AuthenticationRequest) error {
	result := protocol.Socks4RequestGranted
	if auth.Command == protocol.CmdBind {
		result = protocol.Socks4RequestRejected
	}
	socks4Response := protocol.NewSocks4AuthenticationResponse(result, auth.Port, auth.IP[:])

	responseBuffer := alloc.NewSmallBuffer().Clear()
	socks4Response.Write(responseBuffer)
	writer.Write(responseBuffer.Value)
	responseBuffer.Release()

	if result == protocol.Socks4RequestRejected {
		log.Warning("Socks: Unsupported socks 4 command ", auth.Command)
		return ErrorUnsupportedSocksCommand
	}

	dest := v2net.TCPDestination(v2net.IPAddress(auth.IP[:]), auth.Port)
	packet := v2net.NewPacket(dest, nil, true)
	this.transport(reader, writer, packet)
	return nil
}

func (this *SocksServer) transport(reader io.Reader, writer io.Writer, firstPacket v2net.Packet) {
	ray := this.packetDispatcher.DispatchToOutbound(firstPacket)
	input := ray.InboundInput()
	output := ray.InboundOutput()

	var inputFinish, outputFinish sync.Mutex
	inputFinish.Lock()
	outputFinish.Lock()

	go func() {
		v2io.RawReaderToChan(input, reader)
		inputFinish.Unlock()
		close(input)
	}()

	go func() {
		v2io.ChanToRawWriter(writer, output)
		outputFinish.Unlock()
	}()
	outputFinish.Lock()
}
