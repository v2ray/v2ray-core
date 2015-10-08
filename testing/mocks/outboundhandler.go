package mocks

import (
	"bytes"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type OutboundConnectionHandler struct {
	Data2Send   *bytes.Buffer
	Data2Return []byte
	Destination v2net.Destination
}

func (handler *OutboundConnectionHandler) Dispatch(packet v2net.Packet, ray core.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	handler.Destination = packet.Destination()
	if packet.Chunk() != nil {
		handler.Data2Send.Write(packet.Chunk().Value)
	}

	go func() {
		for {
			data, open := <-input
			if !open {
				break
			}
			handler.Data2Send.Write(data.Value)
			data.Release()
		}
		response := alloc.NewBuffer()
		response.Clear()
		response.Append(handler.Data2Return)
		output <- response
		close(output)
	}()

	return nil
}

func (handler *OutboundConnectionHandler) Create(point *core.Point, config interface{}) (core.OutboundConnectionHandler, error) {
	return handler, nil
}
