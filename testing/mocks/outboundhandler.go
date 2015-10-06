package mocks

import (
	"bytes"

	"github.com/v2ray/v2ray-core"
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
		handler.Data2Send.Write(packet.Chunk())
	}

	go func() {
		for {
			data, open := <-input
			if !open {
				break
			}
			handler.Data2Send.Write(data)
		}
		dataCopy := make([]byte, len(handler.Data2Return))
		copy(dataCopy, handler.Data2Return)
		output <- dataCopy
		close(output)
	}()

	return nil
}

func (handler *OutboundConnectionHandler) Create(point *core.Point, config interface{}) (core.OutboundConnectionHandler, error) {
	return handler, nil
}
