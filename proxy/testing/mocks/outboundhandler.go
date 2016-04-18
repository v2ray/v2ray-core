package mocks

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	v2io "github.com/v2ray/v2ray-core/common/io"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type OutboundConnectionHandler struct {
	Destination v2net.Destination
	ConnInput   io.Reader
	ConnOutput  io.Writer
}

func (this *OutboundConnectionHandler) Dispatch(packet v2net.Packet, ray ray.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	this.Destination = packet.Destination()
	if packet.Chunk() != nil {
		this.ConnOutput.Write(packet.Chunk().Value)
		packet.Chunk().Release()
	}

	if packet.MoreChunks() {
		writeFinish := &sync.Mutex{}

		writeFinish.Lock()

		go func() {
			v2io.Pipe(input, v2io.NewAdaptiveWriter(this.ConnOutput))
			writeFinish.Unlock()
			input.Release()
		}()

		writeFinish.Lock()
	}

	v2io.Pipe(v2io.NewAdaptiveReader(this.ConnInput), output)
	output.Close()

	return nil
}

func (this *OutboundConnectionHandler) Create(space app.Space, config interface{}) (proxy.OutboundHandler, error) {
	return this, nil
}
