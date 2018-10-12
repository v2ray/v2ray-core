package commander

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal/done"
	"v2ray.com/core/common/vio"
	"v2ray.com/core/transport/pipe"
)

// OutboundListener is a net.Listener for listening gRPC connections.
type OutboundListener struct {
	buffer chan net.Conn
	done   *done.Instance
}

func (l *OutboundListener) add(conn net.Conn) {
	select {
	case l.buffer <- conn:
	case <-l.done.Wait():
		conn.Close() // nolint: errcheck
	default:
		conn.Close() // nolint: errcheck
	}
}

// Accept implements net.Listener.
func (l *OutboundListener) Accept() (net.Conn, error) {
	select {
	case <-l.done.Wait():
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
			c.Close() // nolint: errcheck
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

// Outbound is a outbound.Handler that handles gRPC connections.
type Outbound struct {
	tag      string
	listener *OutboundListener
	access   sync.RWMutex
	closed   bool
}

// Dispatch implements outbound.Handler.
func (co *Outbound) Dispatch(ctx context.Context, link *vio.Link) {
	co.access.RLock()

	if co.closed {
		pipe.CloseError(link.Reader)
		pipe.CloseError(link.Writer)
		co.access.RUnlock()
		return
	}

	closeSignal := done.New()
	c := net.NewConnection(net.ConnectionInputMulti(link.Writer), net.ConnectionOutputMulti(link.Reader), net.ConnectionOnClose(closeSignal))
	co.listener.add(c)
	co.access.RUnlock()
	<-closeSignal.Wait()
}

// Tag implements outbound.Handler.
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
