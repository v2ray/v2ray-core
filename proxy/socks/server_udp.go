package socks

import (
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet/udp"
)

func (v *Server) listenUDP() error {
	v.udpServer = udp.NewServer(v.packetDispatcher)
	udpHub, err := udp.ListenUDP(v.meta.Address, v.meta.Port, udp.ListenOption{Callback: v.handleUDPPayload})
	if err != nil {
		log.Error("Socks: Failed to listen on udp (", v.meta.Address, ":", v.meta.Port, "): ", err)
		return err
	}
	v.udpMutex.Lock()
	v.udpAddress = v2net.UDPDestination(v.config.GetNetAddress(), v.meta.Port)
	v.udpHub = udpHub
	v.udpMutex.Unlock()
	return nil
}

func (v *Server) handleUDPPayload(payload *buf.Buffer, session *proxy.SessionInfo) {
	defer payload.Release()

	source := session.Source
	log.Info("Socks: Client UDP connection from ", source)
	request, data, err := DecodeUDPPacket(payload.Bytes())

	if err != nil {
		log.Error("Socks|Server: Failed to parse UDP request: ", err)
		return
	}

	if len(data) == 0 {
		return
	}

	log.Info("Socks: Send packet to ", request.Destination(), " with ", len(data), " bytes")
	log.Access(source, request.Destination, log.AccessAccepted, "")

	dataBuf := buf.NewSmall()
	dataBuf.Append(data)
	v.udpServer.Dispatch(&proxy.SessionInfo{Source: source, Destination: request.Destination(), Inbound: v.meta}, dataBuf, func(destination v2net.Destination, payload *buf.Buffer) {
		defer payload.Release()

		log.Info("Socks: Writing back UDP response with ", payload.Len(), " bytes to ", destination)

		udpMessage := EncodeUDPPacket(request, payload.Bytes())
		defer udpMessage.Release()

		v.udpMutex.RLock()
		if !v.accepting {
			v.udpMutex.RUnlock()
			return
		}
		nBytes, err := v.udpHub.WriteTo(udpMessage.Bytes(), destination)
		v.udpMutex.RUnlock()

		if err != nil {
			log.Warning("Socks: failed to write UDP message (", nBytes, " bytes) to ", destination, ": ", err)
		}
	})
}
