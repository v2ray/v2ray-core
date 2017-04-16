package tcp

import (
	"fmt"
	"io"
	"net"

	v2net "v2ray.com/core/common/net"
)

type Server struct {
	Port         v2net.Port
	MsgProcessor func(msg []byte) []byte
	SendFirst    []byte
	Listen       v2net.Address
	accepting    bool
	listener     *net.TCPListener
}

func (server *Server) Start() (v2net.Destination, error) {
	listenerAddr := server.Listen
	if listenerAddr == nil {
		listenerAddr = v2net.LocalHostIP
	}
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   listenerAddr.IP(),
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
	for {
		nBytes, err := conn.Read(request)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Failed to read request:", err)
			}
			break
		}
		response := server.MsgProcessor(request[:nBytes])
		if _, err := conn.Write(response); err != nil {
			fmt.Println("Failed to write response:", err)
			break
		}
	}
	conn.Close()
}

func (v *Server) Close() {
	v.accepting = false
	v.listener.Close()
}
