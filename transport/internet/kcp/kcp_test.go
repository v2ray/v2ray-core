package kcp_test

import (
	"context"
	"crypto/rand"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestDialAndListen(t *testing.T) {
	assert := assert.On(t)

	listerner, err := NewListener(v2net.LocalHostIP, v2net.Port(0), internet.ListenOptions{
		Stream: &internet.StreamConfig{
			Protocol: internet.TransportProtocol_MKCP,
			TransportSettings: []*internet.TransportConfig{
				{
					Protocol: internet.TransportProtocol_MKCP,
					Settings: serial.ToTypedMessage(&Config{}),
				},
			},
		},
	})
	assert.Error(err).IsNil()
	port := v2net.Port(listerner.Addr().(*net.UDPAddr).Port)

	go func() {
		for {
			conn, err := listerner.Accept()
			if err != nil {
				break
			}
			go func() {
				payload := make([]byte, 4096)
				for {
					nBytes, err := conn.Read(payload)
					if err != nil {
						break
					}
					for idx, b := range payload[:nBytes] {
						payload[idx] = b ^ 'c'
					}
					conn.Write(payload[:nBytes])
				}
				conn.Close()
			}()
		}
	}()

	ctx := internet.ContextWithTransportSettings(context.Background(), &Config{})
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		clientConn, err := DialKCP(ctx, v2net.UDPDestination(v2net.LocalHostIP, port))
		assert.Error(err).IsNil()
		wg.Add(1)

		go func() {
			clientSend := make([]byte, 1024*1024)
			rand.Read(clientSend)
			go clientConn.Write(clientSend)

			clientReceived := make([]byte, 1024*1024)
			nBytes, _ := io.ReadFull(clientConn, clientReceived)
			assert.Int(nBytes).Equals(len(clientReceived))
			clientConn.Close()

			clientExpected := make([]byte, 1024*1024)
			for idx, b := range clientSend {
				clientExpected[idx] = b ^ 'c'
			}
			assert.Bytes(clientReceived).Equals(clientExpected)

			wg.Done()
		}()
	}

	wg.Wait()
	for i := 0; i < 60 && listerner.ActiveConnections() > 0; i++ {
		time.Sleep(500 * time.Millisecond)
	}
	assert.Int(listerner.ActiveConnections()).Equals(0)

	listerner.Close()
}
