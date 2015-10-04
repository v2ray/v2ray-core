package freedom

import (
	"net"
	"sync"

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

	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	var readMutex, writeMutex sync.Mutex
	readMutex.Lock()
	writeMutex.Lock()

	if chunk := vconn.packet.Chunk(); chunk != nil {
		conn.Write(chunk)
	}

	if !vconn.packet.MoreChunks() {
		writeMutex.Unlock()
	} else {
		go dumpInput(conn, input, &writeMutex)
	}

	go dumpOutput(conn, output, &readMutex, vconn.packet.Destination().IsUDP())

	go func() {
		writeMutex.Lock()
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.CloseWrite()
		}
		readMutex.Lock()
		conn.Close()
	}()

	return nil
}

func dumpInput(conn net.Conn, input <-chan []byte, finish *sync.Mutex) {
	v2net.ChanToWriter(conn, input)
	finish.Unlock()
}

func dumpOutput(conn net.Conn, output chan<- []byte, finish *sync.Mutex, udp bool) {
	defer finish.Unlock()
	defer close(output)

	response, err := v2net.ReadFrom(conn)
	if len(response) > 0 {
		output <- response
	}
	if err != nil {
		return
	}
	if udp {
		return
	}

	v2net.ReaderToChan(output, conn)
}
