package websocket_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"bytes"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
	. "v2ray.com/core/transport/internet/websocket"
)

func Test_listenWSAndDial(t *testing.T) {
	assert := assert.On(t)
	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13146, internet.ListenOptions{
		Stream: &internet.StreamConfig{
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
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
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
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
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
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
			Network: v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
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
	tlsSettings := &v2tls.Config{
		AllowInsecure: true,
		Certificate: []*v2tls.Certificate{
			{
				Certificate: ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "testing", "tls", "cert.pem"), assert),
				Key:         ReadFile(filepath.Join(os.Getenv("GOPATH"), "src", "v2ray.com", "core", "testing", "tls", "key.pem"), assert),
			},
		},
	}

	listen, err := ListenWS(v2net.DomainAddress("localhost"), 13143, internet.ListenOptions{
		Stream: &internet.StreamConfig{
			SecurityType:     serial.GetMessageType(new(v2tls.Config)),
			SecuritySettings: []*serial.TypedMessage{serial.ToTypedMessage(tlsSettings)},
			Network:          v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
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
			SecurityType:     serial.GetMessageType(new(v2tls.Config)),
			SecuritySettings: []*serial.TypedMessage{serial.ToTypedMessage(tlsSettings)},
			Network:          v2net.Network_WebSocket,
			NetworkSettings: []*internet.NetworkSettings{
				{
					Network: v2net.Network_WebSocket,
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

func ReadFile(file string, assert *assert.Assert) []byte {
	b, err := ioutil.ReadFile(file)
	assert.Error(err).IsNil()
	return b
}
