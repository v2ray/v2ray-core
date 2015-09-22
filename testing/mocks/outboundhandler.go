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

func (handler *OutboundConnectionHandler) Start(ray core.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	go func() {
		for {
			data, open := <-input
			if !open {
				break
			}
			handler.Data2Send.Write(data)
		}
		output <- handler.Data2Return
		close(output)
	}()

	return nil
}

func (handler *OutboundConnectionHandler) Create(point *core.Point, config []byte, packet v2net.Packet) (core.OutboundConnectionHandler, error) {
	handler.Destination = packet.Destination()
	return handler, nil
}
