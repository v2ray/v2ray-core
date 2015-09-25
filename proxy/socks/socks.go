package socks

import (
	"io"
	"net"
	"strconv"
	"sync"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/errors"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
)

// SocksServer is a SOCKS 5 proxy server
type SocksServer struct {
	accepting bool
	vPoint    *core.Point
	config    SocksConfig
}

func NewSocksServer(vp *core.Point, rawConfig []byte) *SocksServer {
	config, err := loadConfig(rawConfig)
	if err != nil {
		log.Error("Unable to load socks config: %v", err)
		panic(errors.NewConfigurationError())
	}
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

	reader := v2net.NewTimeOutReader(4, connection)

	auth, auth4, err := protocol.ReadAuthentication(reader)
	if err != nil && !errors.HasCode(err, 1000) {
		log.Error("Socks failed to read authentication: %v", err)
		return err
	}

	var dest v2net.Destination

	// TODO refactor this part
	if errors.HasCode(err, 1000) {
		result := protocol.Socks4RequestGranted
		if auth4.Command == protocol.CmdBind {
			result = protocol.Socks4RequestRejected
		}
		socks4Response := protocol.NewSocks4AuthenticationResponse(result, auth4.Port, auth4.IP[:])
		connection.Write(socks4Response.ToBytes(nil))

		if result == protocol.Socks4RequestRejected {
			return errors.NewInvalidOperationError("Socks4 command " + strconv.Itoa(int(auth4.Command)))
		}

		dest = v2net.NewTCPDestination(v2net.IPAddress(auth4.IP[:], auth4.Port))
	} else {
		expectedAuthMethod := protocol.AuthNotRequired
		if server.config.AuthMethod == JsonAuthMethodUserPass {
			expectedAuthMethod = protocol.AuthUserPass
		}

		if !auth.HasAuthMethod(expectedAuthMethod) {
			authResponse := protocol.NewAuthenticationResponse(protocol.AuthNoMatchingMethod)
			err = protocol.WriteAuthentication(connection, authResponse)
			if err != nil {
				log.Error("Socks failed to write authentication: %v", err)
				return err
			}
			log.Warning("Socks client doesn't support allowed any auth methods.")
			return errors.NewInvalidOperationError("Unsupported auth methods.")
		}

		authResponse := protocol.NewAuthenticationResponse(expectedAuthMethod)
		err = protocol.WriteAuthentication(connection, authResponse)
		if err != nil {
			log.Error("Socks failed to write authentication: %v", err)
			return err
		}
		if server.config.AuthMethod == JsonAuthMethodUserPass {
			upRequest, err := protocol.ReadUserPassRequest(reader)
			if err != nil {
				log.Error("Socks failed to read username and password: %v", err)
				return err
			}
			status := byte(0)
			if !upRequest.IsValid(server.config.Username, server.config.Password) {
				status = byte(0xFF)
			}
			upResponse := protocol.NewSocks5UserPassResponse(status)
			err = protocol.WriteUserPassResponse(connection, upResponse)
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

		response := protocol.NewSocks5Response()

		if request.Command == protocol.CmdBind || request.Command == protocol.CmdUdpAssociate {
			response := protocol.NewSocks5Response()
			response.Error = protocol.ErrorCommandNotSupported
			err = protocol.WriteResponse(connection, response)
			if err != nil {
				log.Error("Socks failed to write response: %v", err)
				return err
			}
			log.Warning("Unsupported socks command %d", request.Command)
			return errors.NewInvalidOperationError("Socks command " + strconv.Itoa(int(request.Command)))
		}

		response.Error = protocol.ErrorSuccess
		response.Port = request.Port
		response.AddrType = request.AddrType
		switch response.AddrType {
		case protocol.AddrTypeIPv4:
			copy(response.IPv4[:], request.IPv4[:])
		case protocol.AddrTypeIPv6:
			copy(response.IPv6[:], request.IPv6[:])
		case protocol.AddrTypeDomain:
			response.Domain = request.Domain
		}
		err = protocol.WriteResponse(connection, response)
		if err != nil {
			log.Error("Socks failed to write response: %v", err)
			return err
		}

		dest = request.Destination()
	}

	ray := server.vPoint.DispatchToOutbound(v2net.NewTCPPacket(dest))
	input := ray.InboundInput()
	output := ray.InboundOutput()
	var readFinish, writeFinish sync.Mutex
	readFinish.Lock()
	writeFinish.Lock()

	go dumpInput(reader, input, &readFinish)
	go dumpOutput(connection, output, &writeFinish)
	writeFinish.Lock()

	return nil
}

func dumpInput(reader io.Reader, input chan<- []byte, finish *sync.Mutex) {
	v2net.ReaderToChan(input, reader)
	finish.Unlock()
	close(input)
}

func dumpOutput(writer io.Writer, output <-chan []byte, finish *sync.Mutex) {
	v2net.ChanToWriter(writer, output)
	finish.Unlock()
}
