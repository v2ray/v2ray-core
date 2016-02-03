package socks

import (
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
	"github.com/v2ray/v2ray-core/transport/hub"
)

func (this *SocksServer) ListenUDP(port v2net.Port) error {
	this.udpServer = hub.NewUDPServer(this.packetDispatcher)
	udpHub, err := hub.ListenUDP(port, this.handleUDPPayload)
	if err != nil {
		log.Error("Socks: Failed to listen on udp port ", port)
		return err
	}
	this.udpMutex.Lock()
	this.udpAddress = v2net.UDPDestination(this.config.Address, port)
	this.udpHub = udpHub
	this.udpMutex.Unlock()
	return nil
}

func (this *SocksServer) handleUDPPayload(payload *alloc.Buffer, source v2net.Destination) {
	log.Info("Socks: Client UDP connection from ", source)
	request, err := protocol.ReadUDPRequest(payload.Value)
	payload.Release()

	if err != nil {
		log.Error("Socks: Failed to parse UDP request: ", err)
		return
	}
	if request.Data.Len() == 0 {
		request.Data.Release()
		return
	}
	if request.Fragment != 0 {
		log.Warning("Socks: Dropping fragmented UDP packets.")
		// TODO handle fragments
		request.Data.Release()
		return
	}

	udpPacket := v2net.NewPacket(request.Destination(), request.Data, false)
	log.Info("Socks: Send packet to ", udpPacket.Destination(), " with ", request.Data.Len(), " bytes")
	this.udpServer.Dispatch(source, udpPacket, func(packet v2net.Packet) {
		response := &protocol.Socks5UDPRequest{
			Fragment: 0,
			Address:  udpPacket.Destination().Address(),
			Port:     udpPacket.Destination().Port(),
			Data:     packet.Chunk(),
		}
		log.Info("Socks: Writing back UDP response with ", response.Data.Len(), " bytes to ", packet.Destination())

		udpMessage := alloc.NewSmallBuffer().Clear()
		response.Write(udpMessage)

		this.udpMutex.RLock()
		if !this.accepting {
			this.udpMutex.RUnlock()
			return
		}
		nBytes, err := this.udpHub.WriteTo(udpMessage.Value, packet.Destination())
		this.udpMutex.RUnlock()
		udpMessage.Release()
		response.Data.Release()
		if err != nil {
			log.Error("Socks: failed to write UDP message (", nBytes, " bytes) to ", packet.Destination(), ": ", err)
		}
	})
}
