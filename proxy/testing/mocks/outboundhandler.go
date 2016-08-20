package mocks

import (
	"io"
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type OutboundConnectionHandler struct {
	Destination v2net.Destination
	ConnInput   io.Reader
	ConnOutput  io.Writer
}

func (this *OutboundConnectionHandler) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	this.Destination = destination
	this.ConnOutput.Write(payload.Value)
	payload.Release()

	writeFinish := &sync.Mutex{}

	writeFinish.Lock()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(this.ConnOutput)
		defer v2writer.Release()

		v2io.Pipe(input, v2writer)
		writeFinish.Unlock()
		input.Release()
	}()

	writeFinish.Lock()

	v2reader := v2io.NewAdaptiveReader(this.ConnInput)
	defer v2reader.Release()

	v2io.Pipe(v2reader, output)
	output.Close()

	return nil
}

func (this *OutboundConnectionHandler) Create(space app.Space, config interface{}, sendThrough v2net.Address) (proxy.OutboundHandler, error) {
	return this, nil
}
