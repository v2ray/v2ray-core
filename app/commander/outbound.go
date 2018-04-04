package commander

import (
	"context"
	"net"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/transport/ray"
)

type OutboundListener struct {
	buffer chan net.Conn
	done   *signal.Done
}

func (l *OutboundListener) add(conn net.Conn) {
	select {
	case l.buffer <- conn:
	case <-l.done.C():
		conn.Close()
	default:
		conn.Close()
	}
}

// Accept implements net.Listener.
func (l *OutboundListener) Accept() (net.Conn, error) {
	select {
	case <-l.done.C():
		return nil, newError("listen closed")
	case c := <-l.buffer:
		return c, nil
	}
}

// Close implement net.Listener.
func (l *OutboundListener) Close() error {
	common.Must(l.done.Close())
L:
	for {
		select {
		case c := <-l.buffer:
			c.Close()
		default:
			break L
		}
	}
	return nil
}

// Addr implements net.Listener.
func (l *OutboundListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: 0,
	}
}

// Outbound is a core.OutboundHandler that handles gRPC connections.
type Outbound struct {
	tag      string
	listener *OutboundListener
	access   sync.RWMutex
	closed   bool
}

// Dispatch implements core.OutboundHandler.
func (co *Outbound) Dispatch(ctx context.Context, r ray.OutboundRay) {
	co.access.RLock()

	if co.closed {
		r.OutboundInput().CloseError()
		r.OutboundOutput().CloseError()
		co.access.RUnlock()
		return
	}

	closeSignal := signal.NewNotifier()
	c := ray.NewConnection(r.OutboundInput(), r.OutboundOutput(), ray.ConnCloseSignal(closeSignal))
	co.listener.add(c)
	co.access.RUnlock()
	<-closeSignal.Wait()
}

// Tag implements core.OutboundHandler.
func (co *Outbound) Tag() string {
	return co.tag
}

// Start implements common.Runnable.
func (co *Outbound) Start() error {
	co.access.Lock()
	co.closed = false
	co.access.Unlock()
	return nil
}

// Close implements common.Closable.
func (co *Outbound) Close() error {
	co.access.Lock()
	defer co.access.Unlock()

	co.closed = true
	return co.listener.Close()
}
