package udp

import (
	"fmt"
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Server struct {
	Port         uint16
	MsgProcessor func(msg []byte) []byte
}

func (server *Server) Start() (v2net.Address, error) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(server.Port),
		Zone: "",
	})
	if err != nil {
		return nil, err
	}
	go server.handleConnection(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return v2net.IPAddress(localAddr.IP, uint16(localAddr.Port)), nil
}

func (server *Server) handleConnection(conn *net.UDPConn) {
	for {
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
