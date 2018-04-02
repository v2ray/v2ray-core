package udp

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/transport/ray"
)

type ResponseCallback func(payload *buf.Buffer)

type connEntry struct {
	inbound ray.InboundRay
	timer   signal.ActivityUpdater
	cancel  context.CancelFunc
}

type Dispatcher struct {
	sync.RWMutex
	conns      map[net.Destination]*connEntry
	dispatcher core.Dispatcher
}

func NewDispatcher(dispatcher core.Dispatcher) *Dispatcher {
	return &Dispatcher{
		conns:      make(map[net.Destination]*connEntry),
		dispatcher: dispatcher,
	}
}

func (v *Dispatcher) RemoveRay(dest net.Destination) {
	v.Lock()
	defer v.Unlock()
	if conn, found := v.conns[dest]; found {
		conn.inbound.InboundInput().Close()
		conn.inbound.InboundOutput().Close()
		delete(v.conns, dest)
	}
}

func (v *Dispatcher) getInboundRay(dest net.Destination, callback ResponseCallback) *connEntry {
	v.Lock()
	defer v.Unlock()

	if entry, found := v.conns[dest]; found {
		return entry
	}

	newError("establishing new connection for ", dest).WriteToLog()

	ctx, cancel := context.WithCancel(context.Background())
	removeRay := func() {
		cancel()
		v.RemoveRay(dest)
	}
	timer := signal.CancelAfterInactivity(ctx, removeRay, time.Second*4)
	inboundRay, _ := v.dispatcher.Dispatch(ctx, dest)
	entry := &connEntry{
		inbound: inboundRay,
		timer:   timer,
		cancel:  removeRay,
	}
	v.conns[dest] = entry
	go handleInput(ctx, entry, callback)
	return entry
}

func (v *Dispatcher) Dispatch(ctx context.Context, destination net.Destination, payload *buf.Buffer, callback ResponseCallback) {
	// TODO: Add user to destString
	newError("dispatch request to: ", destination).AtDebug().WithContext(ctx).WriteToLog()

	conn := v.getInboundRay(destination, callback)
	outputStream := conn.inbound.InboundInput()
	if outputStream != nil {
		if err := outputStream.WriteMultiBuffer(buf.NewMultiBufferValue(payload)); err != nil {
			newError("failed to write first UDP payload").Base(err).WithContext(ctx).WriteToLog()
			conn.cancel()
			return
		}
	}
}

func handleInput(ctx context.Context, conn *connEntry, callback ResponseCallback) {
	input := conn.inbound.InboundOutput()
	timer := conn.timer

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mb, err := input.ReadMultiBuffer()
		if err != nil {
			newError("failed to handle UDP input").Base(err).WithContext(ctx).WriteToLog()
			conn.cancel()
			return
		}
		timer.Update()
		for _, b := range mb {
			callback(b)
		}
	}
}
