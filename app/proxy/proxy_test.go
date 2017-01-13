package proxy_test

import (
	"context"
	"testing"

	"v2ray.com/core/app"
	. "v2ray.com/core/app/proxy"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/transport/internet"
	_ "v2ray.com/core/transport/internet/tcp"
)

func TestProxyDial(t *testing.T) {
	assert := assert.On(t)

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)
	assert.Error(app.AddApplicationToSpace(ctx, new(proxyman.OutboundConfig))).IsNil()
	outboundManager := proxyman.OutboundHandlerManagerFromSpace(space)
	freedom, err := freedom.New(proxy.ContextWithOutboundMeta(ctx, &proxy.OutboundHandlerMeta{
		Tag: "tag",
		StreamSettings: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_TCP,
		},
	}), &freedom.Config{})
	assert.Error(err).IsNil()
	common.Must(outboundManager.SetHandler("tag", freedom))

	assert.Error(app.AddApplicationToSpace(ctx, new(Config))).IsNil()
	proxy := OutboundProxyFromSpace(space)
	assert.Error(space.Initialize()).IsNil()

	xor := func(b []byte) []byte {
		for idx, x := range b {
			b[idx] = x ^ 'c'
		}
		return b
	}
	tcpServer := &tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()

	conn, err := proxy.Dial(net.LocalHostIP, dest, internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_TCP,
		},
		Proxy: &internet.ProxyConfig{
			Tag: "tag",
		},
	})
	assert.Error(err).IsNil()

	_, err = conn.Write([]byte{'a', 'b', 'c', 'd'})
	assert.Error(err).IsNil()

	b := make([]byte, 10)
	nBytes, err := conn.Read(b)
	assert.Error(err).IsNil()

	assert.Bytes(xor(b[:nBytes])).Equals([]byte{'a', 'b', 'c', 'd'})

	common.Must(conn.Close())
	tcpServer.Close()
}
