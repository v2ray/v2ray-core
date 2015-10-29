package mocks

import (
	"bytes"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

type InboundConnectionHandler struct {
	Data2Send    []byte
	DataReturned *bytes.Buffer
	Port         uint16
	Dispatcher   app.PacketDispatcher
}

func (handler *InboundConnectionHandler) Listen(port uint16) error {
	handler.Port = port
	return nil
}

func (handler *InboundConnectionHandler) Communicate(packet v2net.Packet) error {
	ray := handler.Dispatcher.DispatchToOutbound(packet)

	input := ray.InboundInput()
	output := ray.InboundOutput()

	buffer := alloc.NewBuffer()
	buffer.Clear()
	buffer.Append(handler.Data2Send)
	input <- buffer
	close(input)

	v2net.ChanToWriter(handler.DataReturned, output)
	return nil
}

func (handler *InboundConnectionHandler) Create(dispatcher app.PacketDispatcher, config interface{}) (connhandler.InboundConnectionHandler, error) {
	handler.Dispatcher = dispatcher
	return handler, nil
}
