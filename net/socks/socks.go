package socks

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/v2ray/v2ray-core"
	socksio "github.com/v2ray/v2ray-core/io/socks"
	"github.com/v2ray/v2ray-core/log"
	v2net "github.com/v2ray/v2ray-core/net"
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

	reader := bufio.NewReader(connection)

	auth, auth4, err := socksio.ReadAuthentication(reader)
	if err != nil && err != socksio.ErrorSocksVersion4 {
		log.Error("Error on reading authentication: %v", err)
		return err
	}

	var dest v2net.Address

	// TODO refactor this part
	if err == socksio.ErrorSocksVersion4 {
		result := socksio.Socks4RequestGranted
		if auth4.Command == socksio.CmdBind {
			result = socksio.Socks4RequestRejected
		}
		socks4Response := socksio.NewSocks4AuthenticationResponse(result, auth4.Port, auth4.IP[:])
		socksio.WriteSocks4AuthenticationResponse(connection, socks4Response)

		if result == socksio.Socks4RequestRejected {
			return ErrorCommandNotSupported
		}

		dest = v2net.IPAddress(auth4.IP[:], auth4.Port)
	} else {
		expectedAuthMethod := socksio.AuthNotRequired
		if server.config.AuthMethod == JsonAuthMethodUserPass {
			expectedAuthMethod = socksio.AuthUserPass
		}

		if !auth.HasAuthMethod(expectedAuthMethod) {
			authResponse := socksio.NewAuthenticationResponse(socksio.AuthNoMatchingMethod)
			err = socksio.WriteAuthentication(connection, authResponse)
			if err != nil {
				log.Error("Error on socksio write authentication: %v", err)
				return err
			}
			log.Warning("Client doesn't support allowed any auth methods.")
			return ErrorAuthenticationFailed
		}

		authResponse := socksio.NewAuthenticationResponse(expectedAuthMethod)
		err = socksio.WriteAuthentication(connection, authResponse)
		if err != nil {
			log.Error("Error on socksio write authentication: %v", err)
			return err
		}
		if server.config.AuthMethod == JsonAuthMethodUserPass {
			upRequest, err := socksio.ReadUserPassRequest(reader)
			if err != nil {
				log.Error("Failed to read username and password: %v", err)
				return err
			}
			status := byte(0)
			if !upRequest.IsValid(server.config.Username, server.config.Password) {
				status = byte(0xFF)
			}
			upResponse := socksio.NewSocks5UserPassResponse(status)
			err = socksio.WriteUserPassResponse(connection, upResponse)
			if err != nil {
				log.Error("Error on socksio write user pass response: %v", err)
				return err
			}
			if status != byte(0) {
				return ErrorInvalidUser
			}
		}

		request, err := socksio.ReadRequest(reader)
		if err != nil {
			log.Error("Error on reading socks request: %v", err)
			return err
		}

		response := socksio.NewSocks5Response()

		if request.Command == socksio.CmdBind || request.Command == socksio.CmdUdpAssociate {
			response := socksio.NewSocks5Response()
			response.Error = socksio.ErrorCommandNotSupported
			err = socksio.WriteResponse(connection, response)
			if err != nil {
				log.Error("Error on socksio write response: %v", err)
				return err
			}
			log.Warning("Unsupported socks command %d", request.Command)
			return ErrorCommandNotSupported
		}

		response.Error = socksio.ErrorSuccess
		response.Port = request.Port
		response.AddrType = request.AddrType
		switch response.AddrType {
		case socksio.AddrTypeIPv4:
			copy(response.IPv4[:], request.IPv4[:])
		case socksio.AddrTypeIPv6:
			copy(response.IPv6[:], request.IPv6[:])
		case socksio.AddrTypeDomain:
			response.Domain = request.Domain
		}
		err = socksio.WriteResponse(connection, response)
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
