package socks

import (
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
	"github.com/v2ray/v2ray-core/transport/hub"
)

func (this *Server) listenUDP(address v2net.Address, port v2net.Port) error {
	this.udpServer = hub.NewUDPServer(this.packetDispatcher)
	udpHub, err := hub.ListenUDP(address, port, this.handleUDPPayload)
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

func (this *Server) handleUDPPayload(payload *alloc.Buffer, source v2net.Destination) {
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

	log.Info("Socks: Send packet to ", request.Destination(), " with ", request.Data.Len(), " bytes")
	this.udpServer.Dispatch(source, request.Destination(), request.Data, func(destination v2net.Destination, payload *alloc.Buffer) {
		response := &protocol.Socks5UDPRequest{
			Fragment: 0,
			Address:  request.Destination().Address(),
			Port:     request.Destination().Port(),
			Data:     payload,
		}
		log.Info("Socks: Writing back UDP response with ", payload.Len(), " bytes to ", destination)

		udpMessage := alloc.NewSmallBuffer().Clear()
		response.Write(udpMessage)

		this.udpMutex.RLock()
		if !this.accepting {
			this.udpMutex.RUnlock()
			return
		}
		nBytes, err := this.udpHub.WriteTo(udpMessage.Value, destination)
		this.udpMutex.RUnlock()
		udpMessage.Release()
		response.Data.Release()
		if err != nil {
			log.Error("Socks: failed to write UDP message (", nBytes, " bytes) to ", destination, ": ", err)
		}
	})
}
