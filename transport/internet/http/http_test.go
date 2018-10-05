package http_test

import (
	"context"
	"crypto/rand"
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/tls/cert"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/http"
	"v2ray.com/core/transport/internet/tls"
	. "v2ray.com/ext/assert"
)

func TestHTTPConnection(t *testing.T) {
	assert := With(t)

	port := tcp.PickPort()

	lctx := internet.ContextWithStreamSettings(context.Background(), &internet.MemoryStreamConfig{
		ProtocolName:     "http",
		ProtocolSettings: &Config{},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			Certificate: []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil, cert.CommonName("www.v2ray.com")))},
		},
	})

	listener, err := Listen(lctx, net.LocalHostIP, port, func(conn internet.Connection) {
		go func() {
			defer conn.Close()

			b := buf.New()
			defer b.Release()

			for {
				if err := b.Reset(buf.ReadFrom(conn)); err != nil {
					return
				}
				nBytes, err := conn.Write(b.Bytes())
				assert(err, IsNil)
				assert(int32(nBytes), Equals, b.Len())
			}
		}()
	})
	assert(err, IsNil)

	defer listener.Close()

	time.Sleep(time.Second)

	dctx := internet.ContextWithStreamSettings(context.Background(), &internet.MemoryStreamConfig{
		ProtocolName:     "http",
		ProtocolSettings: &Config{},
		SecurityType:     "tls",
		SecuritySettings: &tls.Config{
			ServerName:    "www.v2ray.com",
			AllowInsecure: true,
		},
	})
	conn, err := Dial(dctx, net.TCPDestination(net.LocalHostIP, port))
	assert(err, IsNil)
	defer conn.Close()

	const N = 1024
	b1 := make([]byte, N)
	common.Must2(rand.Read(b1))
	b2 := buf.New()

	nBytes, err := conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	assert(b2.Reset(buf.ReadFullFrom(conn, N)), IsNil)
	assert(b2.Bytes(), Equals, b1)

	nBytes, err = conn.Write(b1)
	assert(nBytes, Equals, N)
	assert(err, IsNil)

	assert(b2.Reset(buf.ReadFullFrom(conn, N)), IsNil)
	assert(b2.Bytes(), Equals, b1)
}
