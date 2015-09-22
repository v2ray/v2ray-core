package freedom

import (
	"net"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type FreedomConnection struct {
	packet v2net.Packet
}

func NewFreedomConnection(firstPacket v2net.Packet) *FreedomConnection {
	return &FreedomConnection{
		packet: firstPacket,
	}
}

func (vconn *FreedomConnection) Start(ray core.OutboundRay) error {
	conn, err := net.Dial(vconn.packet.Destination().Network(), vconn.packet.Destination().Address().String())
	log.Info("Freedom: Opening connection to %s", vconn.packet.Destination().String())
	if err != nil {
		if ray != nil {
			close(ray.OutboundOutput())
		}
		return log.Error("Freedom: Failed to open connection: %s : %v", vconn.packet.Destination().String(), err)
	}

	if chunk := vconn.packet.Chunk(); chunk != nil {
		conn.Write(chunk)
	}

	if !vconn.packet.MoreChunks() {
		if ray != nil {
			close(ray.OutboundOutput())
		}
		return nil
	}

	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	readFinish := make(chan bool)
	writeFinish := make(chan bool)

	go dumpInput(conn, input, writeFinish)
	go dumpOutput(conn, output, readFinish)

	<-writeFinish
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	<-readFinish
	conn.Close()
	return nil
}

func dumpInput(conn net.Conn, input <-chan []byte, finish chan<- bool) {
	v2net.ChanToWriter(conn, input)
	close(finish)
}

func dumpOutput(conn net.Conn, output chan<- []byte, finish chan<- bool) {
	v2net.ReaderToChan(output, conn)
	close(output)
	close(finish)
}
