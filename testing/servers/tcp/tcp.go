package tcp

import (
	"fmt"
	"net"

	v2io "github.com/v2ray/v2ray-core/common/io"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Server struct {
	Port         v2net.Port
	MsgProcessor func(msg []byte) []byte
	accepting    bool
	listener     *net.TCPListener
}

func (server *Server) Start() (v2net.Destination, error) {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(server.Port),
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
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
			fmt.Printf("Failed accept TCP connection: %v", err)
			continue
		}

		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	for true {
		request, err := v2io.ReadFrom(conn, nil)
		if err != nil {
			break
		}
		response := server.MsgProcessor(request.Value)
		conn.Write(response)
	}
	conn.Close()
}

func (this *Server) Close() {
	this.accepting = false
	this.listener.Close()
}
