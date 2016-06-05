package socks

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
	"github.com/v2ray/v2ray-core/transport/hub"
)

var (
	ErrorUnsupportedSocksCommand = errors.New("Unsupported socks command.")
	ErrorUnsupportedAuthMethod   = errors.New("Unsupported auth method.")
)

// Server is a SOCKS 5 proxy server
type Server struct {
	tcpMutex         sync.RWMutex
	udpMutex         sync.RWMutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *Config
	tcpListener      *hub.TCPHub
	udpHub           *hub.UDPHub
	udpAddress       v2net.Destination
	udpServer        *hub.UDPServer
	meta             *proxy.InboundHandlerMeta
}

// NewServer creates a new Server object.
func NewServer(config *Config, packetDispatcher dispatcher.PacketDispatcher, meta *proxy.InboundHandlerMeta) *Server {
	return &Server{
		config:           config,
		packetDispatcher: packetDispatcher,
		meta:             meta,
	}
}

// Port implements InboundHandler.Port().
func (this *Server) Port() v2net.Port {
	return this.meta.Port
}

// Close implements InboundHandler.Close().
func (this *Server) Close() {
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
func (this *Server) Start() error {
	if this.accepting {
		return nil
	}

	listener, err := hub.ListenTCP(
		this.meta.Address,
		this.meta.Port,
		this.handleConnection,
		nil)
	if err != nil {
		log.Error("Socks: failed to listen on ", this.meta.Address, ":", this.meta.Port, ": ", err)
		return err
	}
	this.accepting = true
	this.tcpMutex.Lock()
	this.tcpListener = listener
	this.tcpMutex.Unlock()
	if this.config.UDPEnabled {
		this.listenUDP()
	}
	return nil
}

func (this *Server) handleConnection(connection *hub.Connection) {
	defer connection.Close()

	timedReader := v2net.NewTimeOutReader(120, connection)
	reader := v2io.NewBufferedReader(timedReader)
	defer reader.Release()

	writer := v2io.NewBufferedWriter(connection)
	defer writer.Release()

	auth, auth4, err := protocol.ReadAuthentication(reader)
	if err != nil && err != protocol.Socks4Downgrade {
		if err != io.EOF {
			log.Warning("Socks: failed to read authentication: ", err)
		}
		return
	}

	cliendAddr := connection.RemoteAddr().String()
	if err != nil && err == protocol.Socks4Downgrade {
		this.handleSocks4(clientAddr, reader, writer, auth4)
	} else {
		this.handleSocks5(clientAddr, reader, writer, auth)
	}
}

func (this *Server) handleSocks5(clientAddr string, reader *v2io.BufferedReader, writer *v2io.BufferedWriter, auth protocol.Socks5AuthenticationRequest) error {
	expectedAuthMethod := protocol.AuthNotRequired
	if this.config.AuthType == AuthTypePassword {
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
		return ErrorUnsupportedAuthMethod
	}

	authResponse := protocol.NewAuthenticationResponse(expectedAuthMethod)
	protocol.WriteAuthentication(writer, authResponse)
	err := writer.Flush()
	if err != nil {
		log.Error("Socks: failed to write authentication: ", err)
		return err
	}
	if this.config.AuthType == AuthTypePassword {
		upRequest, err := protocol.ReadUserPassRequest(reader)
		if err != nil {
			log.Warning("Socks: failed to read username and password: ", err)
			return err
		}
		status := byte(0)
		if !this.config.HasAccount(upRequest.Username(), upRequest.Password()) {
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
			log.Access(clientAddr, "", log.AccessRejected, proxy.ErrorInvalidAuthentication)
			return proxy.ErrorInvalidAuthentication
		}
	}

	request, err := protocol.ReadRequest(reader)
	if err != nil {
		log.Warning("Socks: failed to read request: ", err)
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

		response.Write(writer)
		writer.Flush()
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

	response.Write(writer)
	if err != nil {
		log.Error("Socks: failed to write response: ", err)
		return err
	}

	reader.SetCached(false)
	writer.SetCached(false)

	dest := request.Destination()
	log.Info("Socks: TCP Connect request to ", dest)
	log.Access(clientAddr, dest, log.AccessAccepted, "")

	this.transport(reader, writer, dest)
	return nil
}

func (this *Server) handleUDP(reader io.Reader, writer *v2io.BufferedWriter) error {
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

	response.Write(writer)
	err := writer.Flush()

	if err != nil {
		log.Error("Socks: failed to write response: ", err)
		return err
	}

	// The TCP connection closes after this method returns. We need to wait until
	// the client closes it.
	// TODO: get notified from UDP part
	<-time.After(5 * time.Minute)

	return nil
}

func (this *Server) handleSocks4(clientAddr string, reader *v2io.BufferedReader, writer *v2io.BufferedWriter, auth protocol.Socks4AuthenticationRequest) error {
	result := protocol.Socks4RequestGranted
	if auth.Command == protocol.CmdBind {
		result = protocol.Socks4RequestRejected
	}
	socks4Response := protocol.NewSocks4AuthenticationResponse(result, auth.Port, auth.IP[:])

	socks4Response.Write(writer)

	if result == protocol.Socks4RequestRejected {
		log.Warning("Socks: Unsupported socks 4 command ", auth.Command)
		log.Access(clientAddr, "", log.AccessRejected, ErrorUnsupportedSocksCommand)
		return ErrorUnsupportedSocksCommand
	}

	reader.SetCached(false)
	writer.SetCached(false)

	dest := v2net.TCPDestination(v2net.IPAddress(auth.IP[:]), auth.Port)
	log.Access(clientAddr, dest, log.AccessAccepted, "")
	this.transport(reader, writer, dest)
	return nil
}

func (this *Server) transport(reader io.Reader, writer io.Writer, destination v2net.Destination) {
	ray := this.packetDispatcher.DispatchToOutbound(destination)
	input := ray.InboundInput()
	output := ray.InboundOutput()

	var inputFinish, outputFinish sync.Mutex
	inputFinish.Lock()
	outputFinish.Lock()

	go func() {
		v2reader := v2io.NewAdaptiveReader(reader)
		defer v2reader.Release()

		v2io.Pipe(v2reader, input)
		inputFinish.Unlock()
		input.Close()
	}()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(writer)
		defer v2writer.Release()

		v2io.Pipe(output, v2writer)
		outputFinish.Unlock()
		output.Release()
	}()
	outputFinish.Lock()
}

func init() {
	internal.MustRegisterInboundHandlerCreator("socks",
		func(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
			if !space.HasApp(dispatcher.APP_ID) {
				return nil, internal.ErrorBadConfiguration
			}
			return NewServer(
				rawConfig.(*Config),
				space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher),
				meta), nil
		})
}
