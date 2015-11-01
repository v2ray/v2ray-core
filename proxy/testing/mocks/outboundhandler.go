package mocks

import (
	"io"
	"sync"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
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
			v2net.ChanToWriter(this.ConnOutput, input)
			writeFinish.Unlock()
		}()

		writeFinish.Lock()
	}

	v2net.ReaderToChan(output, this.ConnInput)
	close(output)

	return nil
}

func (this *OutboundConnectionHandler) Create(config interface{}) (connhandler.OutboundConnectionHandler, error) {
	return this, nil
}
