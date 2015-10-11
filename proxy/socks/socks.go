package socks

import (
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/errors"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	jsonconfig "github.com/v2ray/v2ray-core/proxy/socks/config/json"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
)

// SocksServer is a SOCKS 5 proxy server
type SocksServer struct {
	accepting bool
	vPoint    *core.Point
	config    *jsonconfig.SocksConfig
}

func NewSocksServer(vp *core.Point, config *jsonconfig.SocksConfig) *SocksServer {
	return &SocksServer{
		vPoint: vp,
		config: config,
	}
}

func (server *SocksServer) Listen(port uint16) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Error("Socks failed to listen on port %d: %v", port, err)
		return err
	}
	server.accepting = true
	go server.AcceptConnections(listener)
	if server.config.UDPEnabled {
		server.ListenUDP(port)
	}
	return nil
}

func (server *SocksServer) AcceptConnections(listener net.Listener) {
	for server.accepting {
		connection, err := listener.Accept()
		if err != nil {
			log.Error("Socks failed to accept new connection %v", err)
			return
		}
		go server.HandleConnection(connection)
	}
}

func (server *SocksServer) HandleConnection(connection net.Conn) error {
	defer connection.Close()

	reader := v2net.NewTimeOutReader(120, connection)

	auth, auth4, err := protocol.ReadAuthentication(reader)
	if err != nil && !errors.HasCode(err, 1000) {
		log.Error("Socks failed to read authentication: %v", err)
		return err
	}

	if err != nil && errors.HasCode(err, 1000) {
		return server.handleSocks4(reader, connection, auth4)
	} else {
		return server.handleSocks5(reader, connection, auth)
	}
}

func (server *SocksServer) handleSocks5(reader *v2net.TimeOutReader, writer io.Writer, auth protocol.Socks5AuthenticationRequest) error {
	expectedAuthMethod := protocol.AuthNotRequired
	if server.config.IsPassword() {
		expectedAuthMethod = protocol.AuthUserPass
	}

	if !auth.HasAuthMethod(expectedAuthMethod) {
		authResponse := protocol.NewAuthenticationResponse(protocol.AuthNoMatchingMethod)
		err := protocol.WriteAuthentication(writer, authResponse)
		if err != nil {
			log.Error("Socks failed to write authentication: %v", err)
			return err
		}
		log.Warning("Socks client doesn't support allowed any auth methods.")
		return errors.NewInvalidOperationError("Unsupported auth methods.")
	}

	authResponse := protocol.NewAuthenticationResponse(expectedAuthMethod)
	err := protocol.WriteAuthentication(writer, authResponse)
	if err != nil {
		log.Error("Socks failed to write authentication: %v", err)
		return err
	}
	if server.config.IsPassword() {
		upRequest, err := protocol.ReadUserPassRequest(reader)
		if err != nil {
			log.Error("Socks failed to read username and password: %v", err)
			return err
		}
		status := byte(0)
		if !server.config.HasAccount(upRequest.Username(), upRequest.Password()) {
			status = byte(0xFF)
		}
		upResponse := protocol.NewSocks5UserPassResponse(status)
		err = protocol.WriteUserPassResponse(writer, upResponse)
		if err != nil {
			log.Error("Socks failed to write user pass response: %v", err)
			return err
		}
		if status != byte(0) {
			err = errors.NewAuthenticationError(upRequest.AuthDetail())
			log.Warning(err.Error())
			return err
		}
	}

	request, err := protocol.ReadRequest(reader)
	if err != nil {
		log.Error("Socks failed to read request: %v", err)
		return err
	}

	if request.Command == protocol.CmdUdpAssociate && server.config.UDPEnabled {
		return server.handleUDP(reader, writer)
	}

	response := protocol.NewSocks5Response()
	if request.Command == protocol.CmdBind || request.Command == protocol.CmdUdpAssociate {
		response := protocol.NewSocks5Response()
		response.Error = protocol.ErrorCommandNotSupported

		responseBuffer := alloc.NewSmallBuffer().Clear()
		response.Write(responseBuffer)
		_, err = writer.Write(responseBuffer.Value)
		responseBuffer.Release()
		if err != nil {
			log.Error("Socks failed to write response: %v", err)
			return err
		}
		log.Warning("Unsupported socks command %d", request.Command)
		return errors.NewInvalidOperationError("Socks command " + strconv.Itoa(int(request.Command)))
	}

	response.Error = protocol.ErrorSuccess

	// Some SOCKS software requires a value other than dest. Let's fake one:
	response.Port = uint16(1717)
	response.AddrType = protocol.AddrTypeIPv4
	response.IPv4[0] = 0
	response.IPv4[1] = 0
	response.IPv4[2] = 0
	response.IPv4[3] = 0

	responseBuffer := alloc.NewSmallBuffer().Clear()
	response.Write(responseBuffer)
	_, err = writer.Write(responseBuffer.Value)
	responseBuffer.Release()
	if err != nil {
		log.Error("Socks failed to write response: %v", err)
		return err
	}

	dest := request.Destination()
	data, err := v2net.ReadFrom(reader, nil)
	if err != nil {
		return err
	}

	packet := v2net.NewPacket(dest, data, true)
	server.transport(reader, writer, packet)
	return nil
}

func (server *SocksServer) handleUDP(reader *v2net.TimeOutReader, writer io.Writer) error {
	response := protocol.NewSocks5Response()
	response.Error = protocol.ErrorSuccess

	udpAddr := server.getUDPAddr()

	response.Port = udpAddr.Port()
	switch {
	case udpAddr.IsIPv4():
		response.AddrType = protocol.AddrTypeIPv4
		copy(response.IPv4[:], udpAddr.IP())
	case udpAddr.IsIPv6():
		response.AddrType = protocol.AddrTypeIPv6
		copy(response.IPv6[:], udpAddr.IP())
	case udpAddr.IsDomain():
		response.AddrType = protocol.AddrTypeDomain
		response.Domain = udpAddr.Domain()
	}

	responseBuffer := alloc.NewSmallBuffer().Clear()
	response.Write(responseBuffer)
	_, err := writer.Write(responseBuffer.Value)
	responseBuffer.Release()

	if err != nil {
		log.Error("Socks failed to write response: %v", err)
		return err
	}

	reader.SetTimeOut(300)      /* 5 minutes */
	v2net.ReadFrom(reader, nil) // Just in case of anything left in the socket
	// The TCP connection closes after this method returns. We need to wait until
	// the client closes it.
	// TODO: get notified from UDP part
	<-time.After(5 * time.Minute)

	return nil
}

func (server *SocksServer) handleSocks4(reader io.Reader, writer io.Writer, auth protocol.Socks4AuthenticationRequest) error {
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
		return errors.NewInvalidOperationError("Socks4 command " + strconv.Itoa(int(auth.Command)))
	}

	dest := v2net.NewTCPDestination(v2net.IPAddress(auth.IP[:], auth.Port))
	data, err := v2net.ReadFrom(reader, nil)
	if err != nil {
		return err
	}

	packet := v2net.NewPacket(dest, data, true)
	server.transport(reader, writer, packet)
	return nil
}

func (server *SocksServer) transport(reader io.Reader, writer io.Writer, firstPacket v2net.Packet) {
	ray := server.vPoint.DispatchToOutbound(firstPacket)
	input := ray.InboundInput()
	output := ray.InboundOutput()

	var inputFinish, outputFinish sync.Mutex
	inputFinish.Lock()
	outputFinish.Lock()

	go dumpInput(reader, input, &inputFinish)
	go dumpOutput(writer, output, &outputFinish)
	outputFinish.Lock()
}

func dumpInput(reader io.Reader, input chan<- *alloc.Buffer, finish *sync.Mutex) {
	v2net.ReaderToChan(input, reader)
	finish.Unlock()
	close(input)
}

func dumpOutput(writer io.Writer, output <-chan *alloc.Buffer, finish *sync.Mutex) {
	v2net.ChanToWriter(writer, output)
	finish.Unlock()
}
