package udp

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
)

type ResponseCallback func(ctx context.Context, payload *buf.Buffer)

type connEntry struct {
	link   *core.Link
	timer  signal.ActivityUpdater
	cancel context.CancelFunc
}

type Dispatcher struct {
	sync.RWMutex
	conns      map[net.Destination]*connEntry
	dispatcher core.Dispatcher
	callback   ResponseCallback
}

func NewDispatcher(dispatcher core.Dispatcher, callback ResponseCallback) *Dispatcher {
	return &Dispatcher{
		conns:      make(map[net.Destination]*connEntry),
		dispatcher: dispatcher,
		callback:   callback,
	}
}

func (v *Dispatcher) RemoveRay(dest net.Destination) {
	v.Lock()
	defer v.Unlock()
	if conn, found := v.conns[dest]; found {
		common.Close(conn.link.Reader)
		common.Close(conn.link.Writer)
		delete(v.conns, dest)
	}
}

func (v *Dispatcher) getInboundRay(ctx context.Context, dest net.Destination) *connEntry {
	v.Lock()
	defer v.Unlock()

	if entry, found := v.conns[dest]; found {
		return entry
	}

	newError("establishing new connection for ", dest).WriteToLog()

	ctx, cancel := context.WithCancel(ctx)
	removeRay := func() {
		cancel()
		v.RemoveRay(dest)
	}
	timer := signal.CancelAfterInactivity(ctx, removeRay, time.Second*4)
	link, _ := v.dispatcher.Dispatch(ctx, dest)
	entry := &connEntry{
		link:   link,
		timer:  timer,
		cancel: removeRay,
	}
	v.conns[dest] = entry
	go handleInput(ctx, entry, v.callback)
	return entry
}

func (v *Dispatcher) Dispatch(ctx context.Context, destination net.Destination, payload *buf.Buffer) {
	// TODO: Add user to destString
	newError("dispatch request to: ", destination).AtDebug().WriteToLog(session.ExportIDToError(ctx))

	conn := v.getInboundRay(ctx, destination)
	outputStream := conn.link.Writer
	if outputStream != nil {
		if err := outputStream.WriteMultiBuffer(buf.NewMultiBufferValue(payload)); err != nil {
			newError("failed to write first UDP payload").Base(err).WriteToLog(session.ExportIDToError(ctx))
			conn.cancel()
			return
		}
	}
}

func handleInput(ctx context.Context, conn *connEntry, callback ResponseCallback) {
	defer conn.cancel()

	input := conn.link.Reader
	timer := conn.timer

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		mb, err := input.ReadMultiBuffer()
		if err != nil {
			newError("failed to handle UDP input").Base(err).WriteToLog(session.ExportIDToError(ctx))
			return
		}
		timer.Update()
		for _, b := range mb {
			callback(ctx, b)
		}
	}
}
