package mocks

import (
	"bytes"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/net"
)

type FakeOutboundConnectionHandler struct {
	Data2Send   *bytes.Buffer
	Data2Return []byte
	Destination v2net.Address
}

func (handler *FakeOutboundConnectionHandler) Start(ray core.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	output <- handler.Data2Return
	for {
		data, open := <-input
		if !open {
			break
		}
		handler.Data2Send.Write(data)
	}
	return nil
}

func (handler *FakeOutboundConnectionHandler) Create(point *core.Point, dest v2net.Address) (core.OutboundConnectionHandler, error) {
	return handler, nil
}
