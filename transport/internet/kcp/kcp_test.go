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
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestDialAndListen(t *testing.T) {
	assert := assert.On(t)

	conns := make(chan internet.Connection, 16)
	listerner, err := NewListener(internet.ContextWithTransportSettings(context.Background(), &Config{}), v2net.LocalHostIP, v2net.Port(0), conns)
	assert.Error(err).IsNil()
	port := v2net.Port(listerner.Addr().(*net.UDPAddr).Port)

	go func() {
		for conn := range conns {
			go func(c internet.Connection) {
				payload := make([]byte, 4096)
				for {
					nBytes, err := c.Read(payload)
					if err != nil {
						break
					}
					for idx, b := range payload[:nBytes] {
						payload[idx] = b ^ 'c'
					}
					c.Write(payload[:nBytes])
				}
				c.Close()
			}(conn)
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
	close(conns)
}
