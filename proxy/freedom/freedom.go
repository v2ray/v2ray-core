package freedom

import (
	"io"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/transport/dialer"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type FreedomConnection struct {
}

func (this *FreedomConnection) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	log.Info("Freedom: Opening connection to ", destination)

	defer payload.Release()
	defer ray.OutboundInput().Release()
	defer ray.OutboundOutput().Close()

	var conn net.Conn
	err := retry.Timed(5, 100).On(func() error {
		rawConn, err := dialer.Dial(destination)
		if err != nil {
			return err
		}
		conn = rawConn
		return nil
	})
	if err != nil {
		log.Error("Freedom: Failed to open connection to ", destination, ": ", err)
		return err
	}
	defer conn.Close()

	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	var readMutex, writeMutex sync.Mutex
	readMutex.Lock()
	writeMutex.Lock()

	conn.Write(payload.Value)

	go func() {
		v2writer := v2io.NewAdaptiveWriter(conn)
		defer v2writer.Release()

		v2io.Pipe(input, v2writer)
		writeMutex.Unlock()
	}()

	go func() {
		defer readMutex.Unlock()

		var reader io.Reader = conn

		if destination.IsUDP() {
			reader = v2net.NewTimeOutReader(16 /* seconds */, conn)
		}

		v2reader := v2io.NewAdaptiveReader(reader)
		defer v2reader.Release()

		v2io.Pipe(v2reader, output)
		ray.OutboundOutput().Close()
	}()

	writeMutex.Lock()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	readMutex.Lock()

	return nil
}
