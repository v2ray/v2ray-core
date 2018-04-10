package websocket_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/tls/cert"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
	. "v2ray.com/core/transport/internet/websocket"
	. "v2ray.com/ext/assert"
)

func Test_listenWSAndDial(t *testing.T) {
	assert := With(t)
	listen, err := ListenWS(internet.ContextWithTransportSettings(context.Background(), &Config{
		Path: "ws",
	}), net.DomainAddress("localhost"), 13146, func(conn internet.Connection) {
		go func(c internet.Connection) {
			defer c.Close()

			var b [1024]byte
			n, err := c.Read(b[:])
			//assert(err, IsNil)
			if err != nil {
				return
			}
			assert(bytes.HasPrefix(b[:n], []byte("Test connection")), IsTrue)

			_, err = c.Write([]byte("Response"))
			assert(err, IsNil)
		}(conn)
	})
	assert(err, IsNil)

	ctx := internet.ContextWithTransportSettings(context.Background(), &Config{Path: "ws"})
	conn, err := Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13146))

	assert(err, IsNil)
	_, err = conn.Write([]byte("Test connection 1"))
	assert(err, IsNil)

	var b [1024]byte
	n, err := conn.Read(b[:])
	assert(err, IsNil)
	assert(string(b[:n]), Equals, "Response")

	assert(conn.Close(), IsNil)
	<-time.After(time.Second * 5)
	conn, err = Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13146))
	assert(err, IsNil)
	_, err = conn.Write([]byte("Test connection 2"))
	assert(err, IsNil)
	n, err = conn.Read(b[:])
	assert(err, IsNil)
	assert(string(b[:n]), Equals, "Response")
	assert(conn.Close(), IsNil)

	assert(listen.Close(), IsNil)
}

func TestDialWithRemoteAddr(t *testing.T) {
	assert := With(t)
	listen, err := ListenWS(internet.ContextWithTransportSettings(context.Background(), &Config{
		Path: "ws",
	}), net.DomainAddress("localhost"), 13148, func(conn internet.Connection) {
		go func(c internet.Connection) {
			defer c.Close()

			assert(c.RemoteAddr().String(), HasPrefix, "1.1.1.1")

			var b [1024]byte
			n, err := c.Read(b[:])
			//assert(err, IsNil)
			if err != nil {
				return
			}
			assert(bytes.HasPrefix(b[:n], []byte("Test connection")), IsTrue)

			_, err = c.Write([]byte("Response"))
			assert(err, IsNil)
		}(conn)
	})
	assert(err, IsNil)

	ctx := internet.ContextWithTransportSettings(context.Background(), &Config{Path: "ws", Header: []*Header{{Key: "X-Forwarded-For", Value: "1.1.1.1"}}})
	conn, err := Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13148))

	assert(err, IsNil)
	_, err = conn.Write([]byte("Test connection 1"))
	assert(err, IsNil)

	var b [1024]byte
	n, err := conn.Read(b[:])
	assert(err, IsNil)
	assert(string(b[:n]), Equals, "Response")

	assert(listen.Close(), IsNil)
}

func Test_listenWSAndDial_TLS(t *testing.T) {
	assert := With(t)

	start := time.Now()

	ctx := internet.ContextWithTransportSettings(context.Background(), &Config{
		Path: "wss",
	})
	ctx = internet.ContextWithSecuritySettings(ctx, &tls.Config{
		AllowInsecure: true,
		Certificate:   []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil, cert.CommonName("localhost")))},
	})
	listen, err := ListenWS(ctx, net.DomainAddress("localhost"), 13143, func(conn internet.Connection) {
		go func() {
			_ = conn.Close()
		}()
	})
	assert(err, IsNil)
	defer listen.Close()

	conn, err := Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13143))
	assert(err, IsNil)
	_ = conn.Close()

	end := time.Now()
	assert(end.Before(start.Add(time.Second*5)), IsTrue)
}
