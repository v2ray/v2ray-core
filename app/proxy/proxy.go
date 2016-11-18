package proxy

import (
	"errors"
	"io"
	"net"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
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

func (this *OutboundProxy) RegisterDialer() {
	internet.ProxyDialer = this.Dial
}

func (this *OutboundProxy) Dial(src v2net.Address, dest v2net.Destination, options internet.DialerOptions) (internet.Connection, error) {
	handler := this.outboundManager.GetHandler(options.Proxy.Tag)
	if handler == nil {
		log.Warning("Proxy: Failed to get outbound handler with tag: ", options.Proxy.Tag)
		return internet.Dial(src, dest, internet.DialerOptions{
			Stream: options.Stream,
		})
	}
	stream := ray.NewRay()
	go handler.Dispatch(dest, alloc.NewLocalBuffer(32).Clear(), stream)
	return NewProxyConnection(src, dest, stream), nil
}

func (this *OutboundProxy) Release() {

}

type ProxyConnection struct {
	stream     ray.Ray
	closed     bool
	localAddr  net.Addr
	remoteAddr net.Addr

	reader *v2io.ChanReader
	writer *v2io.ChainWriter
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
		reader: v2io.NewChanReader(stream.InboundOutput()),
		writer: v2io.NewChainWriter(stream.InboundInput()),
	}
}

func (this *ProxyConnection) Read(b []byte) (int, error) {
	if this.closed {
		return 0, io.EOF
	}
	return this.reader.Read(b)
}

func (this *ProxyConnection) Write(b []byte) (int, error) {
	if this.closed {
		return 0, io.ErrClosedPipe
	}
	return this.writer.Write(b)
}

func (this *ProxyConnection) Close() error {
	this.closed = true
	this.stream.InboundInput().Close()
	this.stream.InboundOutput().Release()
	return nil
}

func (this *ProxyConnection) LocalAddr() net.Addr {
	return this.localAddr
}

func (this *ProxyConnection) RemoteAddr() net.Addr {
	return this.remoteAddr
}

func (this *ProxyConnection) SetDeadline(t time.Time) error {
	return nil
}

func (this *ProxyConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (this *ProxyConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func (this *ProxyConnection) Reusable() bool {
	return false
}

func (this *ProxyConnection) SetReusable(bool) {

}
