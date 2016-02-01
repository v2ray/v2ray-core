package mocks

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	v2io "github.com/v2ray/v2ray-core/common/io"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type InboundConnectionHandler struct {
	port             v2net.Port
	PacketDispatcher dispatcher.PacketDispatcher
	ConnInput        io.Reader
	ConnOutput       io.Writer
}

func (this *InboundConnectionHandler) Listen(port v2net.Port) error {
	this.port = port
	return nil
}

func (this *InboundConnectionHandler) Port() v2net.Port {
	return this.port
}

func (this *InboundConnectionHandler) Close() {

}

func (this *InboundConnectionHandler) Communicate(packet v2net.Packet) error {
	ray := this.PacketDispatcher.DispatchToOutbound(packet)

	input := ray.InboundInput()
	output := ray.InboundOutput()

	readFinish := &sync.Mutex{}
	writeFinish := &sync.Mutex{}

	readFinish.Lock()
	writeFinish.Lock()

	go func() {
		v2io.RawReaderToChan(input, this.ConnInput)
		close(input)
		readFinish.Unlock()
	}()

	go func() {
		v2io.ChanToRawWriter(this.ConnOutput, output)
		writeFinish.Unlock()
	}()

	readFinish.Lock()
	writeFinish.Lock()
	return nil
}
