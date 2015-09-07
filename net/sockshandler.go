package net

import (
	"net"
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
	go server.AcceptConns(listener)
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

func (server *SocksServer) HandleConnection(connection *net.Conn) error {
	return nil
}
