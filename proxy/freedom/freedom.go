package freedom

import (
	"net"
	"sync"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type FreedomConnection struct {
}

func NewFreedomConnection() *FreedomConnection {
	return &FreedomConnection{}
}

func (vconn *FreedomConnection) Dispatch(firstPacket v2net.Packet, ray core.OutboundRay) error {
	conn, err := net.Dial(firstPacket.Destination().Network(), firstPacket.Destination().Address().String())
	log.Info("Freedom: Opening connection to %s", firstPacket.Destination().String())
	if err != nil {
		close(ray.OutboundOutput())
		log.Error("Freedom: Failed to open connection: %s : %v", firstPacket.Destination().String(), err)
		return err
	}

	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	var readMutex, writeMutex sync.Mutex
	readMutex.Lock()
	writeMutex.Lock()

	if chunk := firstPacket.Chunk(); chunk != nil {
		conn.Write(chunk.Value)
		chunk.Release()
	}

	if !firstPacket.MoreChunks() {
		writeMutex.Unlock()
	} else {
		go dumpInput(conn, input, &writeMutex)
	}

	go dumpOutput(conn, output, &readMutex, firstPacket.Destination().IsUDP())

	writeMutex.Lock()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	readMutex.Lock()
	conn.Close()

	return nil
}

func dumpInput(conn net.Conn, input <-chan *alloc.Buffer, finish *sync.Mutex) {
	v2net.ChanToWriter(conn, input)
	finish.Unlock()
}

func dumpOutput(conn net.Conn, output chan<- *alloc.Buffer, finish *sync.Mutex, udp bool) {
	defer finish.Unlock()
	defer close(output)

	response, err := v2net.ReadFrom(conn, nil)
	log.Info("Freedom receives %d bytes from %s", response.Len(), conn.RemoteAddr().String())
	if response.Len() > 0 {
		output <- response
	} else {
		response.Release()
	}
	if err != nil {
		return
	}
	if udp {
		return
	}

	v2net.ReaderToChan(output, conn)
}
