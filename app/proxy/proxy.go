package proxy

import (
	"io"
	"net"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

const (
	APP_ID = 7
)

type OutboundProxy struct {
	outboundManager proxyman.OutboundHandlerManager
}

func NewOutboundProxy(space app.Space) *OutboundProxy {
	proxy := new(OutboundProxy)
	space.InitializeApplication(func() error {
		if !space.HasApp(proxyman.APP_ID_OUTBOUND_MANAGER) {
			return errors.New("Proxy: Outbound handler manager not found.")
		}
		proxy.outboundManager = space.GetApp(proxyman.APP_ID_OUTBOUND_MANAGER).(proxyman.OutboundHandlerManager)
		return nil
	})
	return proxy
}

func (v *OutboundProxy) RegisterDialer() {
	internet.ProxyDialer = v.Dial
}

// Dial implements internet.Dialer.
func (v *OutboundProxy) Dial(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	handler := v.outboundManager.GetHandler(options.Proxy.Tag)
	if handler == nil {
		log.Warning("Proxy: Failed to get outbound handler with tag: ", options.Proxy.Tag)
		return internet.Dial(src, dest, internet.DialerOptions{
			Stream: options.Stream,
		})
	}
	log.Info("Proxy: Dialing to ", dest)
	stream := ray.NewRay()
	go handler.Dispatch(dest, stream)
	return NewConnection(src, dest, stream), nil
}

type Connection struct {
	stream     ray.Ray
	closed     bool
	localAddr  net.Addr
	remoteAddr net.Addr

	reader *buf.BufferToBytesReader
	writer *buf.BytesToBufferWriter
}

func NewConnection(src v2net.Address, dest v2net.Destination, stream ray.Ray) *Connection {
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
		reader: buf.NewBytesReader(stream.InboundOutput()),
		writer: buf.NewBytesWriter(stream.InboundInput()),
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
	v.stream.InboundOutput().ForceClose()
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

func (v *Connection) SetDeadline(t time.Time) error {
	return nil
}

func (v *Connection) SetReadDeadline(t time.Time) error {
	return nil
}

func (v *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}

func (v *Connection) Reusable() bool {
	return false
}

func (v *Connection) SetReusable(bool) {

}
