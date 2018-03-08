package tcp

import (
	"fmt"
	"io"

	"v2ray.com/core/common/net"
)

type Server struct {
	Port         net.Port
	MsgProcessor func(msg []byte) []byte
	SendFirst    []byte
	Listen       net.Address
	listener     *net.TCPListener
}

func (server *Server) Start() (net.Destination, error) {
	listenerAddr := server.Listen
	if listenerAddr == nil {
		listenerAddr = net.LocalHostIP
	}
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   listenerAddr.IP(),
		Port: int(server.Port),
		Zone: "",
	})
	if err != nil {
		return net.Destination{}, err
	}
	server.Port = net.Port(listener.Addr().(*net.TCPAddr).Port)
	server.listener = listener
	go server.acceptConnections(listener)
	localAddr := listener.Addr().(*net.TCPAddr)
	return net.TCPDestination(net.IPAddress(localAddr.IP), net.Port(localAddr.Port)), nil
}

func (server *Server) acceptConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed accept TCP connection: %v\n", err)
			return
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

func (server *Server) Close() error {
	return server.listener.Close()
}
