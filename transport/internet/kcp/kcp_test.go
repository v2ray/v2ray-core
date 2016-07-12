package kcp_test

import (
	"crypto/rand"
	"io"
	"sync"
	"testing"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/kcp"
)

func TestDialAndListen(t *testing.T) {
	assert := assert.On(t)

	port := v2nettesting.PickPort()
	listerner, err := NewListener(v2net.LocalHostIP, port)
	assert.Error(err).IsNil()

	go func() {
		for {
			conn, err := listerner.Accept()
			if err != nil {
				break
			}
			go func() {
				payload := make([]byte, 1024)
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

	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		clientConn, err := DialKCP(v2net.LocalHostIP, v2net.UDPDestination(v2net.LocalHostIP, port))
		assert.Error(err).IsNil()
		wg.Add(1)

		go func() {
			clientSend := make([]byte, 1024*1024)
			rand.Read(clientSend)
			clientConn.Write(clientSend)

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
	time.Sleep(15 * time.Second)
	assert.Int(listerner.ActiveConnections()).Equals(0)

	listerner.Close()
}
