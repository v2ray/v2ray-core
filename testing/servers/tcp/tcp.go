package tcp

import (
	"fmt"
	"io/ioutil"
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Server struct {
	Port         uint16
	MsgProcessor func(msg []byte) []byte
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
	return v2net.IPAddress(localAddr.IP, uint16(localAddr.Port)), nil
}

func (server *Server) acceptConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed accept TCP connection: %v", err)
			continue
		}

		go server.handleConnection(conn)
	}
}

func (server *Server) handleConnection(conn net.Conn) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Printf("Failed to read request: %v", err)
		return
	}
	response := server.MsgProcessor(request)
	conn.Write(response)
	conn.Close()
}
