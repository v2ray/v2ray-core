package socks

import (
	"errors"
	"net"

	socksio "github.com/v2ray/v2ray-core/io/socks"
)

var (
	ErrorAuthenticationFailed = errors.New("None of the authentication methods is allowed.")
	ErrorCommandNotSupported  = errors.New("Client requested an unsupported command.")
)

// SocksServer is a SOCKS 5 proxy server
type SocksServer struct {
	accepting bool
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

	// TODO: establish connection with VNext

	return nil
}
