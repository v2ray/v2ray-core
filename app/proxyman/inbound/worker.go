package inbound

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal/done"
	"v2ray.com/core/common/task"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tcp"
	"v2ray.com/core/transport/internet/udp"
)

type worker interface {
	Start() error
	Close() error
	Port() net.Port
	Proxy() proxy.Inbound
}

type tcpWorker struct {
	address         net.Address
	port            net.Port
	proxy           proxy.Inbound
	stream          *internet.StreamConfig
	recvOrigDest    bool
	tag             string
	dispatcher      core.Dispatcher
	sniffers        []proxyman.KnownProtocols
	uplinkCounter   core.StatCounter
	downlinkCounter core.StatCounter

	hub internet.Listener
}

func (w *tcpWorker) callback(conn internet.Connection) {
	ctx, cancel := context.WithCancel(context.Background())
	sid := session.NewID()
	ctx = session.ContextWithID(ctx, sid)

	if w.recvOrigDest {
		dest, err := tcp.GetOriginalDestination(conn)
		if err != nil {
			newError("failed to get original destination").Base(err).WriteToLog(session.ExportIDToError(ctx))
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
	if w.uplinkCounter != nil || w.downlinkCounter != nil {
		conn = &internet.StatCouterConnection{
			Connection: conn,
			Uplink:     w.uplinkCounter,
			Downlink:   w.downlinkCounter,
		}
	}
	if err := w.proxy.Process(ctx, net.Network_TCP, conn, w.dispatcher); err != nil {
		newError("connection ends").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}
	cancel()
	if err := conn.Close(); err != nil {
		newError("failed to close connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}
}

func (w *tcpWorker) Proxy() proxy.Inbound {
	return w.proxy
}

func (w *tcpWorker) Start() error {
	ctx := internet.ContextWithStreamSettings(context.Background(), w.stream)
	hub, err := internet.ListenTCP(ctx, w.address, w.port, func(conn internet.Connection) {
		go w.callback(conn)
	})
	if err != nil {
		return newError("failed to listen TCP on ", w.port).AtWarning().Base(err)
	}
	w.hub = hub
	return nil
}

func (w *tcpWorker) Close() error {
	var errors []interface{}
	if w.hub != nil {
		if err := common.Close(w.hub); err != nil {
			errors = append(errors, err)
		}
		if err := common.Close(w.proxy); err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		return newError("failed to close all resources").Base(newError(serial.Concat(errors...)))
	}

	return nil
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
	done             *done.Instance
	uplink           core.StatCounter
	downlink         core.StatCounter
}

func (c *udpConn) updateActivity() {
	atomic.StoreInt64(&c.lastActivityTime, time.Now().Unix())
}

// ReadMultiBuffer implements buf.Reader
func (c *udpConn) ReadMultiBuffer() (buf.MultiBuffer, error) {
	var payload buf.MultiBuffer

	select {
	case in := <-c.input:
		payload.Append(in)
	default:
		select {
		case in := <-c.input:
			payload.Append(in)
		case <-c.done.Wait():
			return nil, io.EOF
		}
	}

L:
	for {
		select {
		case in := <-c.input:
			payload.Append(in)
		default:
			break L
		}
	}

	if c.uplink != nil {
		c.uplink.Add(int64(payload.Len()))
	}

	return payload, nil
}

func (c *udpConn) Read(buf []byte) (int, error) {
	select {
	case in := <-c.input:
		defer in.Release()
		c.updateActivity()
		nBytes := copy(buf, in.Bytes())
		if c.uplink != nil {
			c.uplink.Add(int64(nBytes))
		}
		return nBytes, nil
	case <-c.done.Wait():
		return 0, io.EOF
	}
}

// Write implements io.Writer.
func (c *udpConn) Write(buf []byte) (int, error) {
	n, err := c.output(buf)
	if c.downlink != nil {
		c.downlink.Add(int64(n))
	}
	if err == nil {
		c.updateActivity()
	}
	return n, err
}

func (c *udpConn) Close() error {
	common.Must(c.done.Close())
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

type connID struct {
	src  net.Destination
	dest net.Destination
}

type udpWorker struct {
	sync.RWMutex

	proxy           proxy.Inbound
	hub             *udp.Hub
	address         net.Address
	port            net.Port
	recvOrigDest    bool
	tag             string
	dispatcher      core.Dispatcher
	uplinkCounter   core.StatCounter
	downlinkCounter core.StatCounter

	checker    *task.Periodic
	activeConn map[connID]*udpConn
}

func (w *udpWorker) getConnection(id connID) (*udpConn, bool) {
	w.Lock()
	defer w.Unlock()

	if conn, found := w.activeConn[id]; found && !conn.done.Done() {
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
		done:     done.New(),
		uplink:   w.uplinkCounter,
		downlink: w.downlinkCounter,
	}
	w.activeConn[id] = conn

	conn.updateActivity()
	return conn, false
}

func (w *udpWorker) callback(b *buf.Buffer, source net.Destination, originalDest net.Destination) {
	id := connID{
		src: source,
	}
	if originalDest.IsValid() {
		id.dest = originalDest
	}
	conn, existing := w.getConnection(id)
	select {
	case conn.input <- b:
	case <-conn.done.Wait():
		b.Release()
	default:
		b.Release()
	}

	if !existing {
		go func() {
			ctx := context.Background()
			sid := session.NewID()
			ctx = session.ContextWithID(ctx, sid)

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
			conn.Close() // nolint: errcheck
			w.removeConn(id)
		}()
	}
}

func (w *udpWorker) removeConn(id connID) {
	w.Lock()
	delete(w.activeConn, id)
	w.Unlock()
}

func (w *udpWorker) Start() error {
	w.activeConn = make(map[connID]*udpConn, 16)
	h, err := udp.ListenUDP(w.address, w.port, w.callback, udp.HubReceiveOriginalDestination(w.recvOrigDest), udp.HubCapacity(256))
	if err != nil {
		return err
	}
	w.checker = &task.Periodic{
		Interval: time.Second * 16,
		Execute: func() error {
			nowSec := time.Now().Unix()
			w.Lock()
			defer w.Unlock()

			if len(w.activeConn) == 0 {
				return nil
			}

			for addr, conn := range w.activeConn {
				if nowSec-atomic.LoadInt64(&conn.lastActivityTime) > 8 {
					delete(w.activeConn, addr)
					conn.Close() // nolint: errcheck
				}
			}

			if len(w.activeConn) == 0 {
				w.activeConn = make(map[connID]*udpConn, 16)
			}

			return nil
		},
	}
	if err := w.checker.Start(); err != nil {
		return err
	}
	w.hub = h
	return nil
}

func (w *udpWorker) Close() error {
	w.Lock()
	defer w.Unlock()

	var errors []interface{}

	if w.hub != nil {
		if err := w.hub.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if w.checker != nil {
		if err := w.checker.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if err := common.Close(w.proxy); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return newError("failed to close all resources").Base(newError(serial.Concat(errors...)))
	}
	return nil
}

func (w *udpWorker) Port() net.Port {
	return w.port
}

func (w *udpWorker) Proxy() proxy.Inbound {
	return w.proxy
}
