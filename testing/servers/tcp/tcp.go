package tcp

import (
	"fmt"
	"net"

	v2net "v2ray.com/core/common/net"
)

type Server struct {
	Port         v2net.Port
	MsgProcessor func(msg []byte) []byte
	SendFirst    []byte
	accepting    bool
	listener     *net.TCPListener
}

func (server *Server) Start() (v2net.Destination, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(server.Port),
		Zone: "",
	})
	if err != nil {
		return v2net.Destination{}, err
	}
	server.Port = v2net.Port(listener.Addr().(*net.TCPAddr).Port)
	server.listener = listener
	go server.acceptConnections(listener)
	localAddr := listener.Addr().(*net.TCPAddr)
	return v2net.TCPDestination(v2net.IPAddress(localAddr.IP), v2net.Port(localAddr.Port)), nil
}

func (server *Server) acceptConnections(listener *net.TCPListener) {
	server.accepting = true
	for server.accepting {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed accept TCP connection: %v\n", err)
			continue
		}

		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	if len(server.SendFirst) > 0 {
		conn.Write(server.SendFirst)
	}
	request := make([]byte, 4096)
	for true {
		nBytes, err := conn.Read(request)
		if err != nil {
			break
		}
		response := server.MsgProcessor(request[:nBytes])
		conn.Write(response)
	}
	conn.Close()
}

func (this *Server) Close() {
	this.accepting = false
	this.listener.Close()
}
