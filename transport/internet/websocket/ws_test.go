package websocket_test

import (
	"testing"
	"time"

	"bytes"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
	tlsgen "v2ray.com/core/testing/tls"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
	. "v2ray.com/core/transport/internet/websocket"
)

func Test_listenWSAndDial(t *testing.T) {
	assert := assert.On(t)
	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13146, internet.ListenOptions{
		Stream: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_WebSocket,
			TransportSettings: []*internet.TransportConfig{
				{
					Protocol: internet.TransportProtocol_WebSocket,
					Settings: serial.ToTypedMessage(&Config{
						Path: "ws",
					}),
				},
			},
		},
	})
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
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13146), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_WebSocket,
			TransportSettings: []*internet.TransportConfig{
				{
					Protocol: internet.TransportProtocol_WebSocket,
					Settings: serial.ToTypedMessage(&Config{
						Path: "ws",
					}),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	_, err = conn.Write([]byte("Test connection 1"))
	assert.Error(err).IsNil()

	var b [1024]byte
	n, err := conn.Read(b[:])
	assert.Error(err).IsNil()
	assert.String(string(b[:n])).Equals("Response")

	assert.Error(conn.Close()).IsNil()
	<-time.After(time.Second * 5)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13146), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_WebSocket,
			TransportSettings: []*internet.TransportConfig{
				{
					Protocol: internet.TransportProtocol_WebSocket,
					Settings: serial.ToTypedMessage(&Config{
						Path: "ws",
					}),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	_, err = conn.Write([]byte("Test connection 2"))
	assert.Error(err).IsNil()
	n, err = conn.Read(b[:])
	assert.Error(err).IsNil()
	assert.String(string(b[:n])).Equals("Response")
	assert.Error(conn.Close()).IsNil()
	<-time.After(time.Second * 15)
	conn, err = Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13146), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_WebSocket,
			TransportSettings: []*internet.TransportConfig{
				{
					Protocol: internet.TransportProtocol_WebSocket,
					Settings: serial.ToTypedMessage(&Config{
						Path: "ws",
					}),
				},
			},
		},
	})
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

	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13143, internet.ListenOptions{
		Stream: &internet.StreamConfig{
			SecurityType: serial.GetMessageType(new(v2tls.Config)),
			SecuritySettings: []*serial.TypedMessage{serial.ToTypedMessage(&v2tls.Config{
				Certificate: []*v2tls.Certificate{tlsgen.GenerateCertificateForTest()},
			})},
			Protocol: internet.TransportProtocol_WebSocket,
			TransportSettings: []*internet.TransportConfig{
				{
					Protocol: internet.TransportProtocol_WebSocket,
					Settings: serial.ToTypedMessage(&Config{
						Path: "wss",
						ConnectionReuse: &ConnectionReuse{
							Enable: true,
						},
					}),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	go func() {
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
	}()
	conn, err := Dial(v2net.AnyIP, v2net.TCPDestination(v2net.DomainAddress("localhost"), 13143), internet.DialerOptions{
		Stream: &internet.StreamConfig{
			SecurityType: serial.GetMessageType(new(v2tls.Config)),
			SecuritySettings: []*serial.TypedMessage{serial.ToTypedMessage(&v2tls.Config{
				AllowInsecure: true,
			})},
			Protocol: internet.TransportProtocol_WebSocket,
			TransportSettings: []*internet.TransportConfig{
				{
					Protocol: internet.TransportProtocol_WebSocket,
					Settings: serial.ToTypedMessage(&Config{
						Path: "wss",
						ConnectionReuse: &ConnectionReuse{
							Enable: true,
						},
					}),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	conn.Close()
}
