package proxy_test

import (
	"testing"

	"v2ray.com/core/app"
	. "v2ray.com/core/app/proxy"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
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
	outboundManager := outbound.New()
	common.Must(outboundManager.SetHandler("tag", freedom.New(&freedom.Config{}, space, &proxy.OutboundHandlerMeta{
		Tag: "tag",
		StreamSettings: &internet.StreamConfig{
			Network: v2net.Network_RawTCP,
		},
	})))
	space.BindApp(proxyman.APP_ID_OUTBOUND_MANAGER, outboundManager)

	proxy := NewOutboundProxy(space)
	space.BindApp(APP_ID, proxy)

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

	conn, err := proxy.Dial(v2net.LocalHostIP, dest, internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_RawTCP,
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
