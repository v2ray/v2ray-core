package socks

import (
	_ "bufio"
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	protocol "github.com/v2ray/v2ray-core/proxy/socks/protocol"
)

var (
	ErrorAuthenticationFailed = errors.New("None of the authentication methods is allowed.")
	ErrorCommandNotSupported  = errors.New("Client requested an unsupported command.")
	ErrorInvalidUser          = errors.New("Invalid username or password.")
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
		panic(log.Error("Unable to load socks config: %v", err))
	}
	return &SocksServer{
		vPoint: vp,
		config: config,
	}
}

func (server *SocksServer) Listen(port uint16) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		log.Error("Error on listening port %d: %v", port, err)
		return err
	}
	log.Debug("Working on tcp:%d", port)
	server.accepting = true
	go server.AcceptConnections(listener)
	return nil
}

func (server *SocksServer) AcceptConnections(listener net.Listener) {
	for server.accepting {
		connection, err := listener.Accept()
		if err != nil {
			log.Error("Error on accepting socks connection: %v", err)
		}
		go server.HandleConnection(connection)
	}
}

func (server *SocksServer) HandleConnection(connection net.Conn) error {
	defer connection.Close()

	reader := connection.(io.Reader)

	auth, auth4, err := protocol.ReadAuthentication(reader)
	if err != nil && err != protocol.ErrorSocksVersion4 {
		log.Error("Error on reading authentication: %v", err)
		return err
	}

	var dest *v2net.Destination

	// TODO refactor this part
	if err == protocol.ErrorSocksVersion4 {
		result := protocol.Socks4RequestGranted
		if auth4.Command == protocol.CmdBind {
			result = protocol.Socks4RequestRejected
		}
		socks4Response := protocol.NewSocks4AuthenticationResponse(result, auth4.Port, auth4.IP[:])
		protocol.WriteSocks4AuthenticationResponse(connection, socks4Response)

		if result == protocol.Socks4RequestRejected {
			return ErrorCommandNotSupported
		}

		dest = v2net.NewDestination(v2net.NetTCP, v2net.IPAddress(auth4.IP[:], auth4.Port))
	} else {
		expectedAuthMethod := protocol.AuthNotRequired
		if server.config.AuthMethod == JsonAuthMethodUserPass {
			expectedAuthMethod = protocol.AuthUserPass
		}

		if !auth.HasAuthMethod(expectedAuthMethod) {
			authResponse := protocol.NewAuthenticationResponse(protocol.AuthNoMatchingMethod)
			err = protocol.WriteAuthentication(connection, authResponse)
			if err != nil {
				log.Error("Error on socksio write authentication: %v", err)
				return err
			}
			log.Warning("Client doesn't support allowed any auth methods.")
			return ErrorAuthenticationFailed
		}

		authResponse := protocol.NewAuthenticationResponse(expectedAuthMethod)
		err = protocol.WriteAuthentication(connection, authResponse)
		if err != nil {
			log.Error("Error on socksio write authentication: %v", err)
			return err
		}
		if server.config.AuthMethod == JsonAuthMethodUserPass {
			upRequest, err := protocol.ReadUserPassRequest(reader)
			if err != nil {
				log.Error("Failed to read username and password: %v", err)
				return err
			}
			status := byte(0)
			if !upRequest.IsValid(server.config.Username, server.config.Password) {
				status = byte(0xFF)
			}
			upResponse := protocol.NewSocks5UserPassResponse(status)
			err = protocol.WriteUserPassResponse(connection, upResponse)
			if err != nil {
				log.Error("Error on socksio write user pass response: %v", err)
				return err
			}
			if status != byte(0) {
				return ErrorInvalidUser
			}
		}

		request, err := protocol.ReadRequest(reader)
		if err != nil {
			log.Error("Error on reading socks request: %v", err)
			return err
		}

		response := protocol.NewSocks5Response()

		if request.Command == protocol.CmdBind || request.Command == protocol.CmdUdpAssociate {
			response := protocol.NewSocks5Response()
			response.Error = protocol.ErrorCommandNotSupported
			err = protocol.WriteResponse(connection, response)
			if err != nil {
				log.Error("Error on socksio write response: %v", err)
				return err
			}
			log.Warning("Unsupported socks command %d", request.Command)
			return ErrorCommandNotSupported
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
			log.Error("Error on socksio write response: %v", err)
			return err
		}

		dest = request.Destination()
	}

	ray := server.vPoint.NewInboundConnectionAccepted(dest)
	input := ray.InboundInput()
	output := ray.InboundOutput()
	readFinish := make(chan bool)
	writeFinish := make(chan bool)

	go server.dumpInput(reader, input, readFinish)
	go server.dumpOutput(connection, output, writeFinish)
	<-writeFinish

	return nil
}

func (server *SocksServer) dumpInput(reader io.Reader, input chan<- []byte, finish chan<- bool) {
	v2net.ReaderToChan(input, reader)
	close(input)
	log.Debug("Socks input closed")
	finish <- true
}

func (server *SocksServer) dumpOutput(writer io.Writer, output <-chan []byte, finish chan<- bool) {
	v2net.ChanToWriter(writer, output)
	log.Debug("Socks output closed")
	finish <- true
}
