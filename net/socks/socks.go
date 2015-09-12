package socks

import (
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/v2ray/v2ray-core"
	socksio "github.com/v2ray/v2ray-core/io/socks"
	"github.com/v2ray/v2ray-core/log"
)

var (
	ErrorAuthenticationFailed = errors.New("None of the authentication methods is allowed.")
	ErrorCommandNotSupported  = errors.New("Client requested an unsupported command.")
)

// SocksServer is a SOCKS 5 proxy server
type SocksServer struct {
	accepting bool
	vPoint    *core.VPoint
	config    SocksConfig
}

func NewSocksServer(vp *core.VPoint, rawConfig []byte) *SocksServer {
	server := new(SocksServer)
	server.vPoint = vp
	config, err := loadConfig(rawConfig)
	if err != nil {
		panic(log.Error("Unable to load socks config: %v", err))
	}
	server.config = config
	return server
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

func (server *SocksServer) AcceptConnections(listener net.Listener) error {
	for server.accepting {
		connection, err := listener.Accept()
		if err != nil {
			log.Error("Error on accepting socks connection: %v", err)
			return err
		}
		go server.HandleConnection(connection)
	}
	return nil
}

func (server *SocksServer) HandleConnection(connection net.Conn) error {
	defer connection.Close()

	auth, err := socksio.ReadAuthentication(connection)
	if err != nil {
		log.Error("Error on reading authentication: %v", err)
		return err
	}

	expectedAuthMethod := socksio.AuthNotRequired
	if server.config.AuthMethod == JsonAuthMethodUserPass {
		expectedAuthMethod = socksio.AuthUserPass
	}

	if !auth.HasAuthMethod(expectedAuthMethod) {
		authResponse := socksio.NewAuthenticationResponse(socksio.AuthNoMatchingMethod)
		socksio.WriteAuthentication(connection, authResponse)

		log.Info("Client doesn't support allowed any auth methods.")
		return ErrorAuthenticationFailed
	}

	log.Debug("Auth accepted, responding auth.")
	authResponse := socksio.NewAuthenticationResponse(socksio.AuthNotRequired)
	socksio.WriteAuthentication(connection, authResponse)

	request, err := socksio.ReadRequest(connection)
	if err != nil {
		log.Error("Error on reading socks request: %v", err)
		return err
	}

	response := socksio.NewSocks5Response()

	if request.Command == socksio.CmdBind || request.Command == socksio.CmdUdpAssociate {
		response := socksio.NewSocks5Response()
		response.Error = socksio.ErrorCommandNotSupported
		socksio.WriteResponse(connection, response)
		log.Info("Unsupported socks command %d", request.Command)
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
	log.Debug("Socks response port = %d", response.Port)
	socksio.WriteResponse(connection, response)

	ray := server.vPoint.NewInboundConnectionAccepted(request.Destination())
	input := ray.InboundInput()
	output := ray.InboundOutput()
	finish := make(chan bool, 2)

	go server.dumpInput(connection, input, finish)
	go server.dumpOutput(connection, output, finish)
	server.waitForFinish(finish)

	return nil
}

func (server *SocksServer) dumpInput(conn net.Conn, input chan<- []byte, finish chan<- bool) {
	for {
		buffer := make([]byte, 512)
		nBytes, err := conn.Read(buffer)
		log.Debug("Reading %d bytes, with error %v", nBytes, err)
		if err == io.EOF {
			close(input)
			finish <- true
			log.Debug("Socks finishing input.")
			break
		}
		input <- buffer[:nBytes]
	}
}

func (server *SocksServer) dumpOutput(conn net.Conn, output <-chan []byte, finish chan<- bool) {
	for {
		buffer, open := <-output
		if !open {
			finish <- true
			log.Debug("Socks finishing output")
			break
		}
		nBytes, err := conn.Write(buffer)
		log.Debug("Writing %d bytes with error %v", nBytes, err)
	}
}

func (server *SocksServer) waitForFinish(finish <-chan bool) {
	for i := 0; i < 2; i++ {
		<-finish
	}
}
