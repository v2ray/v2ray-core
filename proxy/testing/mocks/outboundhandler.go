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

func (v *OutboundConnectionHandler) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	v.Destination = destination
	if !payload.IsEmpty() {
		v.ConnOutput.Write(payload.Value)
	}
	payload.Release()

	writeFinish := &sync.Mutex{}

	writeFinish.Lock()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(v.ConnOutput)
		defer v2writer.Release()

		v2io.Pipe(input, v2writer)
		writeFinish.Unlock()
		input.Release()
	}()

	writeFinish.Lock()

	v2reader := v2io.NewAdaptiveReader(v.ConnInput)
	defer v2reader.Release()

	v2io.Pipe(v2reader, output)
	output.Close()
}

func (v *OutboundConnectionHandler) Create(space app.Space, config interface{}, sendThrough v2net.Address) (proxy.OutboundHandler, error) {
	return v, nil
}
