package commander

import (
	"context"
	"net"

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
}

func (co *CommanderOutbound) Dispatch(ctx context.Context, r ray.OutboundRay) {
	closeSignal := signal.NewNotifier()
	c := ray.NewConnection(r.OutboundInput(), r.OutboundOutput(), ray.ConnCloseSignal(closeSignal))
	co.listener.add(c)
	<-closeSignal.Wait()

	return
}

func (co *CommanderOutbound) Tag() string {
	return co.tag
}

func (co *CommanderOutbound) Start() error {
	return nil
}

func (co *CommanderOutbound) Close() {}