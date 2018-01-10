package outbound

import (
	"context"
	"io"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type Handler struct {
	config          *core.OutboundHandlerConfig
	senderSettings  *proxyman.SenderConfig
	proxy           proxy.Outbound
	outboundManager core.OutboundHandlerManager
	mux             *mux.ClientManager
}

func NewHandler(ctx context.Context, config *core.OutboundHandlerConfig) (*Handler, error) {
	v := core.FromContext(ctx)
	if v == nil {
		return nil, newError("V is not in context")
	}
	h := &Handler{
		config:          config,
		outboundManager: v.OutboundHandlerManager(),
	}

	if config.SenderSettings != nil {
		senderSettings, err := config.SenderSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		switch s := senderSettings.(type) {
		case *proxyman.SenderConfig:
			h.senderSettings = s
		default:
			return nil, newError("settings is not SenderConfig")
		}
	}

	proxyConfig, err := config.ProxySettings.GetInstance()
	if err != nil {
		return nil, err
	}

	proxyHandler, err := proxy.CreateOutboundHandler(ctx, proxyConfig)
	if err != nil {
		return nil, err
	}

	if h.senderSettings != nil && h.senderSettings.MultiplexSettings != nil && h.senderSettings.MultiplexSettings.Enabled {
		config := h.senderSettings.MultiplexSettings
		if config.Concurrency < 1 || config.Concurrency > 1024 {
			return nil, newError("invalid mux concurrency: ", config.Concurrency).AtWarning()
		}
		h.mux = mux.NewClientManager(proxyHandler, h, config)
	}

	h.proxy = proxyHandler
	return h, nil
}

func (h *Handler) Tag() string {
	return h.config.Tag
}

// Dispatch implements proxy.Outbound.Dispatch.
func (h *Handler) Dispatch(ctx context.Context, outboundRay ray.OutboundRay) {
	if h.mux != nil {
		err := h.mux.Dispatch(ctx, outboundRay)
		if err != nil {
			newError("failed to process outbound traffic").Base(err).WriteToLog()
			outboundRay.OutboundOutput().CloseError()
		}
	} else {
		err := h.proxy.Process(ctx, outboundRay, h)
		// Ensure outbound ray is properly closed.
		if err != nil && errors.Cause(err) != io.EOF {
			newError("failed to process outbound traffic").Base(err).WriteToLog()
			outboundRay.OutboundOutput().CloseError()
		} else {
			outboundRay.OutboundOutput().Close()
		}
		outboundRay.OutboundInput().CloseError()
	}
}

// Dial implements proxy.Dialer.Dial().
func (h *Handler) Dial(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	if h.senderSettings != nil {
		if h.senderSettings.ProxySettings.HasTag() {
			tag := h.senderSettings.ProxySettings.Tag
			handler := h.outboundManager.GetHandler(tag)
			if handler != nil {
				newError("proxying to ", tag, " for dest ", dest).AtDebug().WriteToLog()
				ctx = proxy.ContextWithTarget(ctx, dest)
				stream := ray.NewRay(ctx)
				go handler.Dispatch(ctx, stream)
				return NewConnection(stream), nil
			}

			newError("failed to get outbound handler with tag: ", tag).AtWarning().WriteToLog()
		}

		if h.senderSettings.Via != nil {
			ctx = internet.ContextWithDialerSource(ctx, h.senderSettings.Via.AsAddress())
		}

		if h.senderSettings.StreamSettings != nil {
			ctx = internet.ContextWithStreamSettings(ctx, h.senderSettings.StreamSettings)
		}
	}

	return internet.Dial(ctx, dest)
}

var (
	_ buf.Reader = (*Connection)(nil)
	_ buf.Writer = (*Connection)(nil)
)

type Connection struct {
	stream     ray.Ray
	closed     bool
	localAddr  net.Addr
	remoteAddr net.Addr

	reader *buf.BufferedReader
	writer buf.Writer
}

func NewConnection(stream ray.Ray) *Connection {
	return &Connection{
		stream: stream,
		localAddr: &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: 0,
		},
		remoteAddr: &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: 0,
		},
		reader: buf.NewBufferedReader(stream.InboundOutput()),
		writer: stream.InboundInput(),
	}
}

// Read implements net.Conn.Read().
func (v *Connection) Read(b []byte) (int, error) {
	if v.closed {
		return 0, io.EOF
	}
	return v.reader.Read(b)
}

func (v *Connection) ReadMultiBuffer() (buf.MultiBuffer, error) {
	return v.reader.ReadMultiBuffer()
}

// Write implements net.Conn.Write().
func (v *Connection) Write(b []byte) (int, error) {
	if v.closed {
		return 0, io.ErrClosedPipe
	}

	l := len(b)
	mb := buf.NewMultiBufferCap(l/buf.Size + 1)
	mb.Write(b)
	return l, v.writer.WriteMultiBuffer(mb)
}

func (v *Connection) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if v.closed {
		return io.ErrClosedPipe
	}

	return v.writer.WriteMultiBuffer(mb)
}

// Close implements net.Conn.Close().
func (v *Connection) Close() error {
	v.closed = true
	v.stream.InboundInput().Close()
	v.stream.InboundOutput().CloseError()
	return nil
}

// LocalAddr implements net.Conn.LocalAddr().
func (v *Connection) LocalAddr() net.Addr {
	return v.localAddr
}

// RemoteAddr implements net.Conn.RemoteAddr().
func (v *Connection) RemoteAddr() net.Addr {
	return v.remoteAddr
}

// SetDeadline implements net.Conn.SetDeadline().
func (v *Connection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements net.Conn.SetReadDeadline().
func (v *Connection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implement net.Conn.SetWriteDeadline().
func (v *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}
