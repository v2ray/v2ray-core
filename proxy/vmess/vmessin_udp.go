package vmess

import (
	"net"

	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
)

const (
	bufferSize = 2 * 1024
)

func (handler *VMessInboundHandler) ListenUDP(port uint16) error {
	addr := &net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Error("VMessIn failed to listen UDP on port %d: %v", port, err)
		return err
	}

	go handler.AcceptPackets(conn)
	return nil
}

func (handler *VMessInboundHandler) AcceptPackets(conn *net.UDPConn) error {
	for {
		buffer := make([]byte, 0, bufferSize)
		nBytes, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Error("VMessIn failed to read UDP packets: %v", err)
			return err
		}
		request, err := protocol.ReadVMessUDP(buffer[:nBytes], handler.clients)
		if err != nil {
			log.Error("VMessIn failed to parse UDP request: %v", err)
			return err
		}

		udpPacket := request.ToPacket()
		go handler.handlePacket(conn, udpPacket, addr)
	}
}

func (handler *VMessInboundHandler) handlePacket(conn *net.UDPConn, packet v2net.Packet, clientAddr *net.UDPAddr) {
	ray := handler.vPoint.DispatchToOutbound(packet)
	close(ray.InboundInput())

	if data, ok := <-ray.InboundOutput(); ok {
		conn.WriteToUDP(data, clientAddr)
	}
}
