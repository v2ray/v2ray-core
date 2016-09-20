package socks

import (
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/socks/protocol"
	"v2ray.com/core/transport/internet/udp"
)

func (this *Server) listenUDP() error {
	this.udpServer = udp.NewUDPServer(this.meta, this.packetDispatcher)
	udpHub, err := udp.ListenUDP(this.meta.Address, this.meta.Port, udp.ListenOption{Callback: this.handleUDPPayload})
	if err != nil {
		log.Error("Socks: Failed to listen on udp ", this.meta.Address, ":", this.meta.Port)
		return err
	}
	this.udpMutex.Lock()
	this.udpAddress = v2net.UDPDestination(this.config.GetNetAddress(), this.meta.Port)
	this.udpHub = udpHub
	this.udpMutex.Unlock()
	return nil
}

func (this *Server) handleUDPPayload(payload *alloc.Buffer, session *proxy.SessionInfo) {
	source := session.Source
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
	log.Access(source, request.Destination, log.AccessAccepted, "")
	this.udpServer.Dispatch(&proxy.SessionInfo{Source: source, Destination: request.Destination()}, request.Data, func(destination v2net.Destination, payload *alloc.Buffer) {
		response := &protocol.Socks5UDPRequest{
			Fragment: 0,
			Address:  request.Destination().Address,
			Port:     request.Destination().Port,
			Data:     payload,
		}
		log.Info("Socks: Writing back UDP response with ", payload.Len(), " bytes to ", destination)

		udpMessage := alloc.NewLocalBuffer(2048).Clear()
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
