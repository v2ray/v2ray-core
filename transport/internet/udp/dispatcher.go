package udp

import (
	"context"
	"sync"

	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/ray"
)

type ResponseCallback func(payload *buf.Buffer)

type Dispatcher struct {
	sync.RWMutex
	conns      map[v2net.Destination]ray.InboundRay
	dispatcher dispatcher.Interface
}

func NewDispatcher(dispatcher dispatcher.Interface) *Dispatcher {
	return &Dispatcher{
		conns:      make(map[v2net.Destination]ray.InboundRay),
		dispatcher: dispatcher,
	}
}

func (v *Dispatcher) RemoveRay(dest v2net.Destination) {
	v.Lock()
	defer v.Unlock()
	if conn, found := v.conns[dest]; found {
		conn.InboundInput().Close()
		conn.InboundOutput().Close()
		delete(v.conns, dest)
	}
}

func (v *Dispatcher) getInboundRay(ctx context.Context, dest v2net.Destination) (ray.InboundRay, bool) {
	v.Lock()
	defer v.Unlock()

	if entry, found := v.conns[dest]; found {
		return entry, true
	}

	log.Trace(errors.New("establishing new connection for ", dest).Path("Transport", "Internet", "UDP", "Dispatcher"))
	inboundRay, _ := v.dispatcher.Dispatch(ctx, dest)
	v.conns[dest] = inboundRay
	return inboundRay, false
}

func (v *Dispatcher) Dispatch(ctx context.Context, destination v2net.Destination, payload *buf.Buffer, callback ResponseCallback) {
	// TODO: Add user to destString
	log.Trace(errors.New("dispatch request to: ", destination).AtDebug().Path("Transport", "Internet", "UDP", "Dispatcher"))

	inboundRay, existing := v.getInboundRay(ctx, destination)
	outputStream := inboundRay.InboundInput()
	if outputStream != nil {
		if err := outputStream.Write(payload); err != nil {
			v.RemoveRay(destination)
		}
	}
	if !existing {
		go func() {
			handleInput(inboundRay.InboundOutput(), callback)
			v.RemoveRay(destination)
		}()
	}
}

func handleInput(input ray.InputStream, callback ResponseCallback) {
	for {
		data, err := input.Read()
		if err != nil {
			break
		}
		callback(data)
	}
}
