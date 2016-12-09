package mocks

import (
	"io"
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/common/buf"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type OutboundConnectionHandler struct {
	Destination v2net.Destination
	ConnInput   io.Reader
	ConnOutput  io.Writer
}

func (v *OutboundConnectionHandler) Dispatch(destination v2net.Destination, payload *buf.Buffer, ray ray.OutboundRay) {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	v.Destination = destination
	if !payload.IsEmpty() {
		v.ConnOutput.Write(payload.Bytes())
	}
	payload.Release()

	writeFinish := &sync.Mutex{}

	writeFinish.Lock()

	go func() {
		v2writer := buf.NewWriter(v.ConnOutput)
		defer v2writer.Release()

		buf.Pipe(input, v2writer)
		writeFinish.Unlock()
		input.Release()
	}()

	writeFinish.Lock()

	v2reader := buf.NewReader(v.ConnInput)
	defer v2reader.Release()

	buf.Pipe(v2reader, output)
	output.Close()
}

func (v *OutboundConnectionHandler) Create(space app.Space, config interface{}, sendThrough v2net.Address) (proxy.OutboundHandler, error) {
	return v, nil
}
