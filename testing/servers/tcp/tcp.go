package tcp

import (
	"fmt"
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Server struct {
	Port         v2net.Port
	MsgProcessor func(msg []byte) []byte
	accepting    bool
}

func (server *Server) Start() (v2net.Address, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(server.Port),
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
	go server.acceptConnections(listener)
	localAddr := listener.Addr().(*net.TCPAddr)
	return v2net.IPAddress(localAddr.IP, v2net.Port(localAddr.Port)), nil
}

func (server *Server) acceptConnections(listener *net.TCPListener) {
	server.accepting = true
	defer listener.Close()
	for server.accepting {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed accept TCP connection: %v", err)
			continue
		}

		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	for true {
		request, err := v2net.ReadFrom(conn, nil)
		if err != nil {
			break
		}
		response := server.MsgProcessor(request.Value)
		conn.Write(response)
	}
	conn.Close()
}

func (this *Server) Close() {
	this.accepting = true
}
