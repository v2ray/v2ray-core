package mocks

import (
	"bytes"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type InboundConnectionHandler struct {
	Data2Send    []byte
	DataReturned *bytes.Buffer
	Port         uint16
	Server       *core.Point
}

func (handler *InboundConnectionHandler) Listen(port uint16) error {
	handler.Port = port
	return nil
}

func (handler *InboundConnectionHandler) Communicate(packet v2net.Packet) error {
	ray := handler.Server.DispatchToOutbound(packet)

	input := ray.InboundInput()
	output := ray.InboundOutput()

	input <- handler.Data2Send
	close(input)

	v2net.ChanToWriter(handler.DataReturned, output)
	return nil
}

func (handler *InboundConnectionHandler) Create(point *core.Point, config []byte) (core.InboundConnectionHandler, error) {
	handler.Server = point
	return handler, nil
}
