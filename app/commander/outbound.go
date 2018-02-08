package commander

import (
	"context"
	"net"
	"sync"

	"v2ray.com/core/common/signal"
	"v2ray.com/core/transport/ray"
)

type OutboundListener struct {
	buffer chan net.Conn
}

func (l *OutboundListener) add(conn net.Conn) {
	select {
	case l.buffer <- conn:
	default:
		conn.Close()
	}
}

func (l *OutboundListener) Accept() (net.Conn, error) {
	c, open := <-l.buffer
	if !open {
		return nil, newError("listener closed")
	}
	return c, nil
}

func (l *OutboundListener) Close() error {
	close(l.buffer)
	return nil
}

func (l *OutboundListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: 0,
	}
}

type CommanderOutbound struct {
	tag      string
	listener *OutboundListener
	access   sync.RWMutex
	closed   bool
}

func (co *CommanderOutbound) Dispatch(ctx context.Context, r ray.OutboundRay) {
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

	return
}

func (co *CommanderOutbound) Tag() string {
	return co.tag
}

func (co *CommanderOutbound) Start() error {
	co.access.Lock()
	co.closed = false
	co.access.Unlock()
	return nil
}

func (co *CommanderOutbound) Close() error {
	co.access.Lock()
	co.closed = true
	co.listener.Close()
	co.access.Unlock()

	return nil
}
