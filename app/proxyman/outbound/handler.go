package outbound

import (
	"context"
	"io"
	"net"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type Handler struct {
	config          *proxyman.OutboundHandlerConfig
	senderSettings  *proxyman.SenderConfig
	proxy           proxy.Outbound
	outboundManager proxyman.OutboundHandlerManager
	mux             *mux.ClientManager
}

func NewHandler(ctx context.Context, config *proxyman.OutboundHandlerConfig) (*Handler, error) {
	h := &Handler{
		config: config,
	}
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("no space in context").Path("App", "Proxyman", "Outbound", "Handler")
	}
	space.OnInitialize(func() error {
		ohm := proxyman.OutboundHandlerManagerFromSpace(space)
		if ohm == nil {
			return errors.New("no OutboundManager in space").Path("App", "Proxyman", "Outbound", "Handler")
		}
		h.outboundManager = ohm
		return nil
	})

	if config.SenderSettings != nil {
		senderSettings, err := config.SenderSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		switch s := senderSettings.(type) {
		case *proxyman.SenderConfig:
			h.senderSettings = s
		default:
			return nil, errors.New("settings is not SenderConfig").Path("App", "Proxyman", "Outbound", "Handler")
		}
	}

	proxyHandler, err := config.GetProxyHandler(ctx)
	if err != nil {
		return nil, err
	}

	if h.senderSettings != nil && h.senderSettings.MultiplexSettings != nil && h.senderSettings.MultiplexSettings.Enabled {
		h.mux = mux.NewClientManager(proxyHandler, h)
	}

	h.proxy = proxyHandler
	return h, nil
}

func (h *Handler) Dispatch(ctx context.Context, outboundRay ray.OutboundRay) {
	if h.mux != nil {
		err := h.mux.Dispatch(ctx, outboundRay)
		if err != nil {
			log.Trace(errors.New("failed to process outbound traffic").Base(err).Path("App", "Proxyman", "Outbound", "Handler"))
		}
	} else {
		err := h.proxy.Process(ctx, outboundRay, h)
		// Ensure outbound ray is properly closed.
		if err != nil && errors.Cause(err) != io.EOF {
			log.Trace(errors.New("failed to process outbound traffic").Base(err).Path("App", "Proxyman", "Outbound", "Handler"))
			outboundRay.OutboundOutput().CloseError()
		} else {
			outboundRay.OutboundOutput().Close()
		}
		outboundRay.OutboundInput().CloseError()
	}
}

// Dial implements proxy.Dialer.Dial().
func (h *Handler) Dial(ctx context.Context, dest v2net.Destination) (internet.Connection, error) {
	if h.senderSettings != nil {
		if h.senderSettings.ProxySettings.HasTag() {
			tag := h.senderSettings.ProxySettings.Tag
			handler := h.outboundManager.GetHandler(tag)
			if handler != nil {
				log.Trace(errors.New("proxying to ", tag).AtDebug().Path("App", "Proxyman", "Outbound", "Handler"))
				ctx = proxy.ContextWithTarget(ctx, dest)
				stream := ray.NewRay(ctx)
				go handler.Dispatch(ctx, stream)
				return NewConnection(stream), nil
			}

			log.Trace(errors.New("failed to get outbound handler with tag: ", tag).AtWarning().Path("App", "Proxyman", "Outbound", "Handler"))
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

type Connection struct {
	stream     ray.Ray
	closed     bool
	localAddr  net.Addr
	remoteAddr net.Addr

	reader io.Reader
	writer io.Writer
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
		reader: buf.ToBytesReader(stream.InboundOutput()),
		writer: buf.ToBytesWriter(stream.InboundInput()),
	}
}

// Read implements net.Conn.Read().
func (v *Connection) Read(b []byte) (int, error) {
	if v.closed {
		return 0, io.EOF
	}
	return v.reader.Read(b)
}

// Write implements net.Conn.Write().
func (v *Connection) Write(b []byte) (int, error) {
	if v.closed {
		return 0, io.ErrClosedPipe
	}
	return v.writer.Write(b)
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

func (v *Connection) Reusable() bool {
	return false
}

func (v *Connection) SetReusable(bool) {

}
