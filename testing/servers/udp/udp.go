package udp

import (
	"fmt"
	"net"

	v2net "v2ray.com/core/common/net"
)

type Server struct {
	Port         v2net.Port
	MsgProcessor func(msg []byte) []byte
	accepting    bool
	conn         *net.UDPConn
}

func (server *Server) Start() (v2net.Destination, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: int(server.Port),
		Zone: "",
	})
	if err != nil {
		return v2net.Destination{}, err
	}
	server.Port = v2net.Port(conn.LocalAddr().(*net.UDPAddr).Port)
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
		if _, err := conn.WriteToUDP(response, addr); err != nil {
			fmt.Println("Failed to write to UDP: ", err.Error())
		}
	}
}

func (server *Server) Close() {
	server.accepting = false
	server.conn.Close()
}
