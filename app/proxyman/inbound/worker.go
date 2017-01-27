package inbound

import (
	"context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core/common/buf"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/udp"
)

type worker interface {
	Start() error
	Close()
	Port() v2net.Port
	Proxy() proxy.Inbound
}

type tcpWorker struct {
	address          v2net.Address
	port             v2net.Port
	proxy            proxy.Inbound
	stream           *internet.StreamConfig
	recvOrigDest     bool
	tag              string
	allowPassiveConn bool

	ctx    context.Context
	cancel context.CancelFunc
	hub    *internet.TCPHub
}

func (w *tcpWorker) callback(conn internet.Connection) {
	ctx, cancel := context.WithCancel(w.ctx)
	if w.recvOrigDest {
		dest := tcp.GetOriginalDestination(conn)
		if dest.IsValid() {
			ctx = proxy.ContextWithOriginalDestination(ctx, dest)
		}
	}
	if len(w.tag) > 0 {
		ctx = proxy.ContextWithInboundTag(ctx, w.tag)
	}
	ctx = proxy.ContextWithAllowPassiveConnection(ctx, w.allowPassiveConn)
	ctx = proxy.ContextWithInboundDestination(ctx, v2net.TCPDestination(w.address, w.port))
	w.proxy.Process(ctx, v2net.Network_TCP, conn)
	cancel()
	conn.Close()
}

func (w *tcpWorker) Proxy() proxy.Inbound {
	return w.proxy
}

func (w *tcpWorker) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	w.ctx = ctx
	w.cancel = cancel
	hub, err := internet.ListenTCP(w.address, w.port, w.callback, w.stream)
	if err != nil {
		return err
	}
	w.hub = hub
	return nil
}

func (w *tcpWorker) Close() {
	w.hub.Close()
	w.cancel()
}

func (w *tcpWorker) Port() v2net.Port {
	return w.port
}

type udpConn struct {
	lastActivityTime int64 // in seconds
	input            chan []byte
	output           func([]byte) (int, error)
	remote           net.Addr
	local            net.Addr
	cancel           context.CancelFunc
}

func (c *udpConn) updateActivity() {
	atomic.StoreInt64(&c.lastActivityTime, time.Now().Unix())
}

func (c *udpConn) Read(buf []byte) (int, error) {
	in, open := <-c.input
	if !open {
		return 0, io.EOF
	}
	c.updateActivity()
	return copy(buf, in), nil
}

func (c *udpConn) Write(buf []byte) (int, error) {
	n, err := c.output(buf)
	if err == nil {
		c.updateActivity()
	}
	return n, err
}

func (c *udpConn) Close() error {
	return nil
}

func (c *udpConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *udpConn) LocalAddr() net.Addr {
	return c.remote
}

func (*udpConn) SetDeadline(time.Time) error {
	return nil
}

func (*udpConn) SetReadDeadline(time.Time) error {
	return nil
}

func (*udpConn) SetWriteDeadline(time.Time) error {
	return nil
}

func (*udpConn) Reusable() bool {
	return false
}

func (*udpConn) SetReusable(bool) {}

type udpWorker struct {
	sync.RWMutex

	proxy        proxy.Inbound
	hub          *udp.Hub
	address      v2net.Address
	port         v2net.Port
	recvOrigDest bool
	tag          string

	ctx        context.Context
	cancel     context.CancelFunc
	activeConn map[v2net.Destination]*udpConn
}

func (w *udpWorker) getConnection(src v2net.Destination) (*udpConn, bool) {
	w.Lock()
	defer w.Unlock()

	if conn, found := w.activeConn[src]; found {
		return conn, true
	}

	conn := &udpConn{
		input: make(chan []byte, 32),
		output: func(b []byte) (int, error) {
			return w.hub.WriteTo(b, src)
		},
		remote: &net.UDPAddr{
			IP:   src.Address.IP(),
			Port: int(src.Port),
		},
		local: &net.UDPAddr{
			IP:   w.address.IP(),
			Port: int(w.port),
		},
	}
	w.activeConn[src] = conn

	conn.updateActivity()
	return conn, false
}

func (w *udpWorker) callback(b *buf.Buffer, source v2net.Destination, originalDest v2net.Destination) {
	conn, existing := w.getConnection(source)
	conn.input <- b.Bytes()

	if !existing {
		go func() {
			ctx := w.ctx
			ctx, cancel := context.WithCancel(ctx)
			conn.cancel = cancel
			if originalDest.IsValid() {
				ctx = proxy.ContextWithOriginalDestination(ctx, originalDest)
			}
			if len(w.tag) > 0 {
				ctx = proxy.ContextWithInboundTag(ctx, w.tag)
			}
			ctx = proxy.ContextWithSource(ctx, source)
			ctx = proxy.ContextWithInboundDestination(ctx, v2net.UDPDestination(w.address, w.port))
			w.proxy.Process(ctx, v2net.Network_UDP, conn)
			w.removeConn(source)
			cancel()
		}()
	}
}

func (w *udpWorker) removeConn(src v2net.Destination) {
	w.Lock()
	delete(w.activeConn, src)
	w.Unlock()
}

func (w *udpWorker) Start() error {
	w.activeConn = make(map[v2net.Destination]*udpConn)
	ctx, cancel := context.WithCancel(context.Background())
	w.ctx = ctx
	w.cancel = cancel
	h, err := udp.ListenUDP(w.address, w.port, udp.ListenOption{
		Callback:            w.callback,
		ReceiveOriginalDest: w.recvOrigDest,
	})
	if err != nil {
		return err
	}
	go w.monitor()
	w.hub = h
	return nil
}

func (w *udpWorker) Close() {
	w.hub.Close()
	w.cancel()
}

func (w *udpWorker) monitor() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case <-time.After(time.Second * 16):
			nowSec := time.Now().Unix()
			w.Lock()
			for addr, conn := range w.activeConn {
				if nowSec-conn.lastActivityTime > 8 {
					delete(w.activeConn, addr)
					conn.cancel()
				}
			}
			w.Unlock()
		}
	}
}

func (w *udpWorker) Port() v2net.Port {
	return w.port
}

func (w *udpWorker) Proxy() proxy.Inbound {
	return w.proxy
}
