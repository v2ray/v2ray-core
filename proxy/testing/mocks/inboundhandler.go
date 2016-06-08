package mocks

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	v2io "github.com/v2ray/v2ray-core/common/io"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type InboundConnectionHandler struct {
	ListeningPort    v2net.Port
	ListeningAddress v2net.Address
	PacketDispatcher dispatcher.PacketDispatcher
	ConnInput        io.Reader
	ConnOutput       io.Writer
}

func (this *InboundConnectionHandler) Start() error {
	return nil
}

func (this *InboundConnectionHandler) Port() v2net.Port {
	return this.ListeningPort
}

func (this *InboundConnectionHandler) Close() {

}

func (this *InboundConnectionHandler) Communicate(destination v2net.Destination) error {
	ray := this.PacketDispatcher.DispatchToOutbound(destination)

	input := ray.InboundInput()
	output := ray.InboundOutput()

	readFinish := &sync.Mutex{}
	writeFinish := &sync.Mutex{}

	readFinish.Lock()
	writeFinish.Lock()

	go func() {
		v2reader := v2io.NewAdaptiveReader(this.ConnInput)
		defer v2reader.Release()

		v2io.Pipe(v2reader, input)
		input.Close()
		readFinish.Unlock()
	}()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(this.ConnOutput)
		defer v2writer.Release()

		v2io.Pipe(output, v2writer)
		output.Release()
		writeFinish.Unlock()
	}()

	readFinish.Lock()
	writeFinish.Lock()
	return nil
}
