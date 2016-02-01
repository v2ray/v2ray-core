package freedom

import (
	"net"
	"sync"

	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/transport/dialer"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type FreedomConnection struct {
}

func (this *FreedomConnection) Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error {
	log.Info("Freedom: Opening connection to ", firstPacket.Destination())

	var conn net.Conn
	err := retry.Timed(5, 100).On(func() error {
		rawConn, err := dialer.Dial(firstPacket.Destination())
		if err != nil {
			return err
		}
		conn = rawConn
		return nil
	})
	if err != nil {
		close(ray.OutboundOutput())
		log.Error("Freedom: Failed to open connection to ", firstPacket.Destination(), ": ", err)
		return err
	}
	defer conn.Close()

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
		go func() {
			v2io.ChanToRawWriter(conn, input)
			writeMutex.Unlock()
		}()
	}

	go func() {
		defer readMutex.Unlock()
		defer close(output)

		response, err := v2io.ReadFrom(conn, nil)
		log.Info("Freedom receives ", response.Len(), " bytes from ", conn.RemoteAddr())
		if response.Len() > 0 {
			output <- response
		} else {
			response.Release()
		}
		if err != nil {
			return
		}
		if firstPacket.Destination().IsUDP() {
			return
		}

		v2io.RawReaderToChan(output, conn)
	}()

	writeMutex.Lock()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	readMutex.Lock()

	return nil
}
