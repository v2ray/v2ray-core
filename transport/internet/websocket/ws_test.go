package websocket_test

import (
	"bytes"
	"context"
	"runtime"
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/tls/cert"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
	. "v2ray.com/core/transport/internet/websocket"
	. "v2ray.com/ext/assert"
)

func Test_listenWSAndDial(t *testing.T) {
	assert := With(t)

	listen, err := ListenWS(context.Background(), net.LocalHostIP, 13146, &internet.MemoryStreamConfig{
		ProtocolName: "websocket",
		ProtocolSettings: &Config{
			Path: "ws",
		},
	}, func(conn internet.Connection) {
		go func(c internet.Connection) {
			defer c.Close()

			var b [1024]byte
			n, err := c.Read(b[:])
			//common.Must(err)
			if err != nil {
				return
			}
			assert(bytes.HasPrefix(b[:n], []byte("Test connection")), IsTrue)

			_, err = c.Write([]byte("Response"))
			common.Must(err)
		}(conn)
	})
	common.Must(err)

	ctx := context.Background()
	streamSettings := &internet.MemoryStreamConfig{
		ProtocolName:     "websocket",
		ProtocolSettings: &Config{Path: "ws"},
	}
	conn, err := Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13146), streamSettings)

	common.Must(err)
	_, err = conn.Write([]byte("Test connection 1"))
	common.Must(err)

	var b [1024]byte
	n, err := conn.Read(b[:])
	common.Must(err)
	assert(string(b[:n]), Equals, "Response")

	assert(conn.Close(), IsNil)
	<-time.After(time.Second * 5)
	conn, err = Dial(ctx, net.TCPDestination(net.DomainAddress("localhost"), 13146), streamSettings)
	common.Must(err)
	_, err = conn.Write([]byte("Test connection 2"))
	common.Must(err)
	n, err = conn.Read(b[:])
	common.Must(err)
	assert(string(b[:n]), Equals, "Response")
	assert(conn.Close(), IsNil)

	assert(listen.Close(), IsNil)
}

func TestDialWithRemoteAddr(t *testing.T) {
	assert := With(t)
	listen, err := ListenWS(context.Background(), net.LocalHostIP, 13148, &internet.MemoryStreamConfig{
		ProtocolName: "websocket",
		ProtocolSettings: &Config{
			Path: "ws",
		},
	}, func(conn internet.Connection) {
		go func(c internet.Connection) {
			defer c.Close()

			assert(c.RemoteAddr().String(), HasPrefix, "1.1.1.1")

			var b [1024]byte
			n, err := c.Read(b[:])
			//common.Must(err)
			if err != nil {
				return
			}
			assert(bytes.HasPrefix(b[:n], []byte("Test connection")), IsTrue)

			_, err = c.Write([]byte("Response"))
			common.Must(err)
		}(conn)
	})
	common.Must(err)

	conn, err := Dial(context.Background(), net.TCPDestination(net.DomainAddress("localhost"), 13148), &internet.MemoryStreamConfig{
		ProtocolName:     "websocket",
		ProtocolSettings: &Config{Path: "ws", Header: []*Header{{Key: "X-Forwarded-For", Value: "1.1.1.1"}}},
	})

	common.Must(err)
	_, err = conn.Write([]byte("Test connection 1"))
	common.Must(err)

	var b [1024]byte
	n, err := conn.Read(b[:])
	common.Must(err)
	assert(string(b[:n]), Equals, "Response")

	assert(listen.Close(), IsNil)
}

func Test_listenWSAndDial_TLS(t *testing.T) {
	if runtime.GOARCH == "arm64" {
		return
	}

	assert := With(t)

	start := time.Now()

	streamSettings := &internet.MemoryStreamConfig{
		ProtocolName: "websocket",
		ProtocolSettings: &Config{
			Path: "wss",
		},
		SecurityType: "tls",
		SecuritySettings: &tls.Config{
			AllowInsecure: true,
			Certificate:   []*tls.Certificate{tls.ParseCertificate(cert.MustGenerate(nil, cert.CommonName("localhost")))},
		},
	}
	listen, err := ListenWS(context.Background(), net.LocalHostIP, 13143, streamSettings, func(conn internet.Connection) {
		go func() {
			_ = conn.Close()
		}()
	})
	common.Must(err)
	defer listen.Close()

	conn, err := Dial(context.Background(), net.TCPDestination(net.DomainAddress("localhost"), 13143), streamSettings)
	common.Must(err)
	_ = conn.Close()

	end := time.Now()
	assert(end.Before(start.Add(time.Second*5)), IsTrue)
}
