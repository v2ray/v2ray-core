package vmess

import (
	"bytes"
	"crypto/md5"
	"net"

	v2io "github.com/v2ray/v2ray-core/common/io"
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
		buffer := make([]byte, bufferSize)
		nBytes, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Error("VMessIn failed to read UDP packets: %v", err)
			continue
		}

		reader := bytes.NewReader(buffer[:nBytes])
		requestReader := protocol.NewVMessRequestReader(handler.clients)

		request, err := requestReader.Read(reader)
		if err != nil {
			log.Warning("VMessIn: Invalid request from (%s): %v", addr.String(), err)
			continue
		}

		cryptReader, err := v2io.NewAesDecryptReader(request.RequestKey[:], request.RequestIV[:], reader)
		if err != nil {
			log.Error("VMessIn: Failed to create decrypt reader: %v", err)
			continue
		}

		data := make([]byte, bufferSize)
		nBytes, err = cryptReader.Read(data)
		if err != nil {
			log.Warning("VMessIn: Unable to decrypt data: %v", err)
			continue
		}

		packet := v2net.NewPacket(request.Destination(), data[:nBytes], false)
		go handler.handlePacket(conn, request, packet, addr)
	}
}

func (handler *VMessInboundHandler) handlePacket(conn *net.UDPConn, request *protocol.VMessRequest, packet v2net.Packet, clientAddr *net.UDPAddr) {
	ray := handler.vPoint.DispatchToOutbound(packet)
	close(ray.InboundInput())

	responseKey := md5.Sum(request.RequestKey[:])
	responseIV := md5.Sum(request.RequestIV[:])

	buffer := bytes.NewBuffer(make([]byte, 0, bufferSize))

	response := protocol.NewVMessResponse(request)
	responseWriter, err := v2io.NewAesEncryptWriter(responseKey[:], responseIV[:], buffer)
	if err != nil {
		log.Error("VMessIn: Failed to create encrypt writer: %v", err)
		return
	}
	responseWriter.Write(response[:])

	hasData := false

	if data, ok := <-ray.InboundOutput(); ok {
		hasData = true
		responseWriter.Write(data)
	}

	if hasData {
		conn.WriteToUDP(buffer.Bytes(), clientAddr)
	}
}
