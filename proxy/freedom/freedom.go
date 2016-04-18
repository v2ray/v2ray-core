package freedom

import (
	"io"
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

	defer firstPacket.Release()
	defer ray.OutboundInput().Release()

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
	}

	if !firstPacket.MoreChunks() {
		writeMutex.Unlock()
	} else {
		go func() {
			v2io.Pipe(input, v2io.NewAdaptiveWriter(conn))
			writeMutex.Unlock()
		}()
	}

	go func() {
		defer readMutex.Unlock()
		defer output.Close()

		var reader io.Reader = conn

		if firstPacket.Destination().IsUDP() {
			reader = v2net.NewTimeOutReader(16 /* seconds */, conn)
		}

		v2io.Pipe(v2io.NewAdaptiveReader(reader), output)
	}()

	writeMutex.Lock()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	readMutex.Lock()

	return nil
}
