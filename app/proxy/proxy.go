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
	go handler.Dispatch(dest, nil, stream)
	return NewProxyConnection(src, dest, stream), nil
}

func (v *OutboundProxy) Release() {

}

type ProxyConnection struct {
	stream     ray.Ray
	closed     bool
	localAddr  net.Addr
	remoteAddr net.Addr

	reader *buf.BufferToBytesReader
	writer *buf.BytesToBufferWriter
}

func NewProxyConnection(src v2net.Address, dest v2net.Destination, stream ray.Ray) *ProxyConnection {
	return &ProxyConnection{
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

func (v *ProxyConnection) Read(b []byte) (int, error) {
	if v.closed {
		return 0, io.EOF
	}
	return v.reader.Read(b)
}

func (v *ProxyConnection) Write(b []byte) (int, error) {
	if v.closed {
		return 0, io.ErrClosedPipe
	}
	return v.writer.Write(b)
}

func (v *ProxyConnection) Close() error {
	v.closed = true
	v.stream.InboundInput().Close()
	v.stream.InboundOutput().Release()
	v.reader.Release()
	v.writer.Release()
	return nil
}

func (v *ProxyConnection) LocalAddr() net.Addr {
	return v.localAddr
}

func (v *ProxyConnection) RemoteAddr() net.Addr {
	return v.remoteAddr
}

func (v *ProxyConnection) SetDeadline(t time.Time) error {
	return nil
}

func (v *ProxyConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (v *ProxyConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func (v *ProxyConnection) Reusable() bool {
	return false
}

func (v *ProxyConnection) SetReusable(bool) {

}
