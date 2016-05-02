package udp

import (
	"fmt"
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Server struct {
	Port         v2net.Port
	MsgProcessor func(msg []byte) []byte
	accepting    bool
	conn         *net.UDPConn
}

func (server *Server) Start() (v2net.Destination, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(server.Port),
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
	server.conn = conn
	go server.handleConnection(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return v2net.UDPDestination(v2net.IPAddress(localAddr.IP), v2net.Port(localAddr.Port)), nil
}

func (server *Server) handleConnection(conn *net.UDPConn) {
	server.accepting = true
	for server.accepting {
		buffer := make([]byte, 2*1024)
		nBytes, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Failed to read from UDP: %v\n", err)
			continue
		}

		response := server.MsgProcessor(buffer[:nBytes])
		conn.WriteToUDP(response, addr)
	}
}

func (server *Server) Close() {
	server.accepting = false
	server.conn.Close()
}
