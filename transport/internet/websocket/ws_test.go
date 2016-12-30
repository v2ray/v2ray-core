package websocket_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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
		conn, err := listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		conn, err = listen.Accept()
		assert.Error(err).IsNil()
		conn.Close()
		listen.Close()
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
	conn.Close()
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
	conn.Close()
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
	conn.Close()
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
