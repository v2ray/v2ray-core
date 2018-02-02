package inbound

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/udp"
)

type worker interface {
	Start() error
	Close()
	Port() net.Port
	Proxy() proxy.Inbound
}

type tcpWorker struct {
	address      net.Address
	port         net.Port
	proxy        proxy.Inbound
	stream       *internet.StreamConfig
	recvOrigDest bool
	tag          string
	dispatcher   core.Dispatcher
	sniffers     []proxyman.KnownProtocols

	ctx    context.Context
	cancel context.CancelFunc
	hub    internet.Listener
}

func (w *tcpWorker) callback(conn internet.Connection) {
	ctx, cancel := context.WithCancel(w.ctx)
	if w.recvOrigDest {
		dest, err := tcp.GetOriginalDestination(conn)
		if err != nil {
			newError("failed to get original destination").Base(err).WriteToLog()
		}
		if dest.IsValid() {
			ctx = proxy.ContextWithOriginalTarget(ctx, dest)
		}
	}
	if len(w.tag) > 0 {
		ctx = proxy.ContextWithInboundTag(ctx, w.tag)
	}
	ctx = proxy.ContextWithInboundEntryPoint(ctx, net.TCPDestination(w.address, w.port))
	ctx = proxy.ContextWithSource(ctx, net.DestinationFromAddr(conn.RemoteAddr()))
	if len(w.sniffers) > 0 {
		ctx = proxyman.ContextWithProtocolSniffers(ctx, w.sniffers)
	}
	if err := w.proxy.Process(ctx, net.Network_TCP, conn, w.dispatcher); err != nil {
		newError("connection ends").Base(err).WriteToLog()
	}
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
	ctx = internet.ContextWithStreamSettings(ctx, w.stream)
	conns := make(chan internet.Connection, 16)
	hub, err := internet.ListenTCP(ctx, w.address, w.port, conns)
	if err != nil {
		return newError("failed to listen TCP on ", w.port).AtWarning().Base(err)
	}
	go w.handleConnections(conns)
	w.hub = hub
	return nil
}

func (w *tcpWorker) handleConnections(conns <-chan internet.Connection) {
	for {
		select {
		case <-w.ctx.Done():
			w.hub.Close()
		L:
			for {
				select {
				case conn := <-conns:
					conn.Close()
				default:
					break L
				}
			}
			return
		case conn := <-conns:
			go w.callback(conn)
		}
	}
}

func (w *tcpWorker) Close() {
	if w.hub != nil {
		w.cancel()
	}
}

func (w *tcpWorker) Port() net.Port {
	return w.port
}

type udpConn struct {
	lastActivityTime int64 // in seconds
	input            chan *buf.Buffer
	output           func([]byte) (int, error)
	remote           net.Addr
	local            net.Addr
	ctx              context.Context
	cancel           context.CancelFunc
}

func (c *udpConn) updateActivity() {
	atomic.StoreInt64(&c.lastActivityTime, time.Now().Unix())
}

func (c *udpConn) Read(buf []byte) (int, error) {
	select {
	case in := <-c.input:
		defer in.Release()
		c.updateActivity()
		return copy(buf, in.Bytes()), nil
	case <-c.ctx.Done():
		return 0, io.EOF
	}
}

// Write implements io.Writer.
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

type connId struct {
	src  net.Destination
	dest net.Destination
}

type udpWorker struct {
	sync.RWMutex

	proxy        proxy.Inbound
	hub          *udp.Hub
	address      net.Address
	port         net.Port
	recvOrigDest bool
	tag          string
	dispatcher   core.Dispatcher

	ctx        context.Context
	cancel     context.CancelFunc
	activeConn map[connId]*udpConn
}

func (w *udpWorker) getConnection(id connId) (*udpConn, bool) {
	w.Lock()
	defer w.Unlock()

	if conn, found := w.activeConn[id]; found {
		return conn, true
	}

	conn := &udpConn{
		input: make(chan *buf.Buffer, 32),
		output: func(b []byte) (int, error) {
			return w.hub.WriteTo(b, id.src)
		},
		remote: &net.UDPAddr{
			IP:   id.src.Address.IP(),
			Port: int(id.src.Port),
		},
		local: &net.UDPAddr{
			IP:   w.address.IP(),
			Port: int(w.port),
		},
	}
	w.activeConn[id] = conn

	conn.updateActivity()
	return conn, false
}

func (w *udpWorker) callback(b *buf.Buffer, source net.Destination, originalDest net.Destination) {
	id := connId{
		src:  source,
		dest: originalDest,
	}
	conn, existing := w.getConnection(id)
	select {
	case conn.input <- b:
	default:
		b.Release()
	}

	if !existing {
		go func() {
			ctx := w.ctx
			ctx, cancel := context.WithCancel(ctx)
			conn.ctx = ctx
			conn.cancel = cancel
			if originalDest.IsValid() {
				ctx = proxy.ContextWithOriginalTarget(ctx, originalDest)
			}
			if len(w.tag) > 0 {
				ctx = proxy.ContextWithInboundTag(ctx, w.tag)
			}
			ctx = proxy.ContextWithSource(ctx, source)
			ctx = proxy.ContextWithInboundEntryPoint(ctx, net.UDPDestination(w.address, w.port))
			if err := w.proxy.Process(ctx, net.Network_UDP, conn, w.dispatcher); err != nil {
				newError("connection ends").Base(err).WriteToLog()
			}
			w.removeConn(id)
			cancel()
		}()
	}
}

func (w *udpWorker) removeConn(id connId) {
	w.Lock()
	delete(w.activeConn, id)
	w.Unlock()
}

func (w *udpWorker) Start() error {
	w.activeConn = make(map[connId]*udpConn, 16)
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
	if w.hub != nil {
		w.hub.Close()
		w.cancel()
	}
}

func (w *udpWorker) monitor() {
	timer := time.NewTicker(time.Second * 16)
	defer timer.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-timer.C:
			nowSec := time.Now().Unix()
			w.Lock()
			for addr, conn := range w.activeConn {
				if nowSec-atomic.LoadInt64(&conn.lastActivityTime) > 8 {
					delete(w.activeConn, addr)
					conn.cancel()
				}
			}
			w.Unlock()
		}
	}
}

func (w *udpWorker) Port() net.Port {
	return w.port
}

func (w *udpWorker) Proxy() proxy.Inbound {
	return w.proxy
}
