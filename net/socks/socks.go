package socks

import (
	"errors"
	"io"
	"net"

	"github.com/v2ray/v2ray-core"
	socksio "github.com/v2ray/v2ray-core/io/socks"
)

var (
	ErrorAuthenticationFailed = errors.New("None of the authentication methods is allowed.")
	ErrorCommandNotSupported  = errors.New("Client requested an unsupported command.")
)

// SocksServer is a SOCKS 5 proxy server
type SocksServer struct {
	accepting bool
	vPoint    *core.VPoint
}

func NewSocksServer(vp *core.VPoint) *SocksServer {
	server := new(SocksServer)
	server.vPoint = vp
	return server
}

func (server *SocksServer) Listen(port uint8) error {
	listener, err := net.Listen("tcp", ":"+string(port))
	if err != nil {
		return err
	}
	server.accepting = true
	go server.AcceptConnections(listener)
	return nil
}

func (server *SocksServer) AcceptConnections(listener net.Listener) error {
	for server.accepting {
		connection, err := listener.Accept()
		if err != nil {
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
		return err
	}

	if auth.HasAuthMethod(socksio.AuthNotRequired) {
		return ErrorAuthenticationFailed
	}

	authResponse := socksio.NewAuthenticationResponse(socksio.AuthNotRequired)
	socksio.WriteAuthentication(connection, authResponse)

	request, err := socksio.ReadRequest(connection)
	if err != nil {
		return err
	}

	if request.Command == socksio.CmdBind || request.Command == socksio.CmdUdpAssociate {
		response := socksio.NewSocks5Response()
		response.Error = socksio.ErrorCommandNotSupported
		socksio.WriteResponse(connection, response)
		return ErrorCommandNotSupported
	}

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
		buffer := make([]byte, 256)
		nBytes, err := conn.Read(buffer)
		if err == io.EOF {
			finish <- true
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
			break
		}
		conn.Write(buffer)
	}
}

func (server *SocksServer) waitForFinish(finish <-chan bool) {
	for i := 0; i < 2; i++ {
		<-finish
	}
}
