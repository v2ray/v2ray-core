package websocket_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	tlsgen "v2ray.com/core/testing/tls"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
	. "v2ray.com/core/transport/internet/websocket"
)

func Test_listenWSAndDial(t *testing.T) {
	assert := assert.On(t)
	listen, err := ListenWS(internet.ContextWithTransportSettings(context.Background(), &Config{
		Path: "ws",
	}), v2net.DomainAddress("localhost"), 13146)
	assert.Error(err).IsNil()
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				break
			}
			go func() {
				defer conn.Close()

				var b [1024]byte
				n, err := conn.Read(b[:])
				//assert.Error(err).IsNil()
				if err != nil {
					conn.SetReusable(false)
					return
				}
				assert.Bool(bytes.HasPrefix(b[:n], []byte("Test connection"))).IsTrue()

				_, err = conn.Write([]byte("Response"))
				assert.Error(err).IsNil()
			}()
		}
	}()

	ctx := internet.ContextWithTransportSettings(context.Background(), &Config{Path: "ws"})
	conn, err := Dial(ctx, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13146))

	assert.Error(err).IsNil()
	_, err = conn.Write([]byte("Test connection 1"))
	assert.Error(err).IsNil()

	var b [1024]byte
	n, err := conn.Read(b[:])
	assert.Error(err).IsNil()
	assert.String(string(b[:n])).Equals("Response")

	assert.Error(conn.Close()).IsNil()
	<-time.After(time.Second * 5)
	conn, err = Dial(ctx, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13146))
	assert.Error(err).IsNil()
	_, err = conn.Write([]byte("Test connection 2"))
	assert.Error(err).IsNil()
	n, err = conn.Read(b[:])
	assert.Error(err).IsNil()
	assert.String(string(b[:n])).Equals("Response")
	assert.Error(conn.Close()).IsNil()
	<-time.After(time.Second * 15)
	conn, err = Dial(ctx, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13146))
	assert.Error(err).IsNil()
	_, err = conn.Write([]byte("Test connection 3"))
	assert.Error(err).IsNil()
	n, err = conn.Read(b[:])
	assert.Error(err).IsNil()
	assert.String(string(b[:n])).Equals("Response")
	assert.Error(conn.Close()).IsNil()

	assert.Error(listen.Close()).IsNil()
}

func Test_listenWSAndDial_TLS(t *testing.T) {
	assert := assert.On(t)
	go func() {
		<-time.After(time.Second * 5)
		assert.Fail("Too slow")
	}()

	ctx := internet.ContextWithTransportSettings(context.Background(), &Config{
		Path: "wss",
		ConnectionReuse: &ConnectionReuse{
			Enable: true,
		},
	})
	ctx = internet.ContextWithSecuritySettings(ctx, &v2tls.Config{
		AllowInsecure: true,
		Certificate:   []*v2tls.Certificate{tlsgen.GenerateCertificateForTest()},
	})
	listen, err := ListenWS(ctx, v2net.DomainAddress("localhost"), 13143)
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()

	conn, err := Dial(ctx, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13143))
	assert.Error(err).IsNil()
	conn.Close()
}
