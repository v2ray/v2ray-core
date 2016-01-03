package socks

import (
	"net"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
)

func (this *SocksServer) ListenUDP(port v2net.Port) error {
	addr := &net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Error("Socks failed to listen UDP on port %d: %v", port, err)
		return err
	}
	this.udpAddress = v2net.UDPDestination(v2net.IPAddress(this.config.IP()), port)
	this.udpConn = conn

	go this.AcceptPackets()
	return nil
}

func (this *SocksServer) AcceptPackets() error {
	for this.accepting {
		buffer := alloc.NewBuffer()
		this.RLock()
		if !this.accepting {
			this.RUnlock()
			return nil
		}
		nBytes, addr, err := this.udpConn.ReadFromUDP(buffer.Value)
		this.RUnlock()
		if err != nil {
			log.Error("Socks failed to read UDP packets: %v", err)
			buffer.Release()
			continue
		}
		log.Info("Client UDP connection from %v", addr)
		request, err := protocol.ReadUDPRequest(buffer.Value[:nBytes])
		buffer.Release()
		if err != nil {
			log.Error("Socks failed to parse UDP request: %v", err)
			continue
		}
		if request.Data == nil || request.Data.Len() == 0 {
			continue
		}
		if request.Fragment != 0 {
			log.Warning("Dropping fragmented UDP packets.")
			// TODO handle fragments
			request.Data.Release()
			continue
		}

		udpPacket := v2net.NewPacket(request.Destination(), request.Data, false)
		log.Info("Send packet to %s with %d bytes", udpPacket.Destination().String(), request.Data.Len())
		go this.handlePacket(udpPacket, addr, request.Address, request.Port)
	}
	return nil
}

func (this *SocksServer) handlePacket(packet v2net.Packet, clientAddr *net.UDPAddr, targetAddr v2net.Address, port v2net.Port) {
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
	close(ray.InboundInput())

	for data := range ray.InboundOutput() {
		response := &protocol.Socks5UDPRequest{
			Fragment: 0,
			Address:  targetAddr,
			Port:     port,
			Data:     data,
		}
		log.Info("Writing back UDP response with %d bytes from %s to %s", data.Len(), targetAddr.String(), clientAddr.String())

		udpMessage := alloc.NewSmallBuffer().Clear()
		response.Write(udpMessage)

		this.RLock()
		if !this.accepting {
			this.RUnlock()
			return
		}
		nBytes, err := this.udpConn.WriteToUDP(udpMessage.Value, clientAddr)
		this.RUnlock()
		udpMessage.Release()
		response.Data.Release()
		if err != nil {
			log.Error("Socks failed to write UDP message (%d bytes) to %s: %v", nBytes, clientAddr.String(), err)
		}
	}
}
