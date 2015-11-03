package vmess

import (
	"bytes"
	"crypto/md5"
	"net"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
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

func (handler *VMessInboundHandler) AcceptPackets(conn *net.UDPConn) {
	for {
		buffer := alloc.NewBuffer()
		nBytes, addr, err := conn.ReadFromUDP(buffer.Value)
		if err != nil {
			log.Error("VMessIn failed to read UDP packets: %v", err)
			buffer.Release()
			continue
		}

		reader := bytes.NewReader(buffer.Value[:nBytes])
		requestReader := protocol.NewVMessRequestReader(handler.clients)

		request, err := requestReader.Read(reader)
		if err != nil {
			log.Access(addr.String(), "", log.AccessRejected, err.Error())
			log.Warning("VMessIn: Invalid request from (%s): %v", addr.String(), err)
			buffer.Release()
			continue
		}
		log.Access(addr.String(), request.Address.String(), log.AccessAccepted, "")

		aesStream, err := v2crypto.NewAesDecryptionStream(request.RequestKey, request.RequestIV)
		if err != nil {
			log.Error("VMessIn: Failed to AES decryption stream: %v", err)
			buffer.Release()
			continue
		}
		cryptReader := v2crypto.NewCryptionReader(aesStream, reader)

		data := alloc.NewBuffer()
		nBytes, err = cryptReader.Read(data.Value)
		buffer.Release()
		if err != nil {
			log.Warning("VMessIn: Unable to decrypt data: %v", err)
			data.Release()
			continue
		}
		data.Slice(0, nBytes)

		packet := v2net.NewPacket(request.Destination(), data, false)
		go handler.handlePacket(conn, request, packet, addr)
	}
}

func (handler *VMessInboundHandler) handlePacket(conn *net.UDPConn, request *protocol.VMessRequest, packet v2net.Packet, clientAddr *net.UDPAddr) {
	ray := handler.dispatcher.DispatchToOutbound(packet)
	close(ray.InboundInput())

	responseKey := md5.Sum(request.RequestKey)
	responseIV := md5.Sum(request.RequestIV)

	buffer := alloc.NewBuffer().Clear()
	defer buffer.Release()

	aesStream, err := v2crypto.NewAesEncryptionStream(responseKey[:], responseIV[:])
	if err != nil {
		log.Error("VMessIn: Failed to create AES encryption stream: %v", err)
		return
	}
	responseWriter := v2crypto.NewCryptionWriter(aesStream, buffer)

	responseWriter.Write(request.ResponseHeader)

	hasData := false

	if data, ok := <-ray.InboundOutput(); ok {
		hasData = true
		responseWriter.Write(data.Value)
		data.Release()
	}

	if hasData {
		conn.WriteToUDP(buffer.Value, clientAddr)
		log.Info("VMessIn sending %d bytes to %s", buffer.Len(), clientAddr.String())
	}
}
