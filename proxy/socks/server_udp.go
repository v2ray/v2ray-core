package socks

import (
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/socks/protocol"
	"v2ray.com/core/transport/internet/udp"
)

func (v *Server) listenUDP() error {
	v.udpServer = udp.NewUDPServer(v.packetDispatcher)
	udpHub, err := udp.ListenUDP(v.meta.Address, v.meta.Port, udp.ListenOption{Callback: v.handleUDPPayload})
	if err != nil {
		log.Error("Socks: Failed to listen on udp ", v.meta.Address, ":", v.meta.Port)
		return err
	}
	v.udpMutex.Lock()
	v.udpAddress = v2net.UDPDestination(v.config.GetNetAddress(), v.meta.Port)
	v.udpHub = udpHub
	v.udpMutex.Unlock()
	return nil
}

func (v *Server) handleUDPPayload(payload *buf.Buffer, session *proxy.SessionInfo) {
	source := session.Source
	log.Info("Socks: Client UDP connection from ", source)
	request, err := protocol.ReadUDPRequest(payload.Bytes())
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
	v.udpServer.Dispatch(&proxy.SessionInfo{Source: source, Destination: request.Destination(), Inbound: v.meta}, request.Data, func(destination v2net.Destination, payload *buf.Buffer) {
		response := &protocol.Socks5UDPRequest{
			Fragment: 0,
			Address:  request.Destination().Address,
			Port:     request.Destination().Port,
			Data:     payload,
		}
		log.Info("Socks: Writing back UDP response with ", payload.Len(), " bytes to ", destination)

		udpMessage := buf.NewLocal(2048)
		response.Write(udpMessage)

		v.udpMutex.RLock()
		if !v.accepting {
			v.udpMutex.RUnlock()
			return
		}
		nBytes, err := v.udpHub.WriteTo(udpMessage.Bytes(), destination)
		v.udpMutex.RUnlock()
		udpMessage.Release()
		response.Data.Release()
		if err != nil {
			log.Error("Socks: failed to write UDP message (", nBytes, " bytes) to ", destination, ": ", err)
		}
	})
}
