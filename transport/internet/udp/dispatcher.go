package udp

import (
	"context"
	"sync"

	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/app/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type ResponseCallback func(payload *buf.Buffer)

type Dispatcher struct {
	sync.RWMutex
	conns            map[string]ray.InboundRay
	packetDispatcher dispatcher.Interface
}

func NewDispatcher(packetDispatcher dispatcher.Interface) *Dispatcher {
	return &Dispatcher{
		conns:            make(map[string]ray.InboundRay),
		packetDispatcher: packetDispatcher,
	}
}

func (v *Dispatcher) RemoveRay(name string) {
	v.Lock()
	defer v.Unlock()
	if conn, found := v.conns[name]; found {
		conn.InboundInput().Close()
		conn.InboundOutput().Close()
		delete(v.conns, name)
	}
}

func (v *Dispatcher) getInboundRay(ctx context.Context, dest v2net.Destination) (ray.InboundRay, bool) {
	destString := dest.String()
	v.Lock()
	defer v.Unlock()

	if entry, found := v.conns[destString]; found {
		return entry, true
	}

	log.Info("UDP|Server: establishing new connection for ", dest)
	ctx = proxy.ContextWithDestination(ctx, dest)
	return v.packetDispatcher.DispatchToOutbound(ctx), false
}

func (v *Dispatcher) Dispatch(ctx context.Context, destination v2net.Destination, payload *buf.Buffer, callback ResponseCallback) {
	// TODO: Add user to destString
	destString := destination.String()
	log.Debug("UDP|Server: Dispatch request: ", destString)

	inboundRay, existing := v.getInboundRay(ctx, destination)
	outputStream := inboundRay.InboundInput()
	if outputStream != nil {
		if err := outputStream.Write(payload); err != nil {
			v.RemoveRay(destString)
		}
	}
	if !existing {
		go func() {
			handleInput(inboundRay.InboundOutput(), callback)
			v.RemoveRay(destString)
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
