package kcp_test

import (
	"context"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestDialAndListen(t *testing.T) {
	listerner, err := NewListener(context.Background(), net.LocalHostIP, net.Port(0), &internet.MemoryStreamConfig{
		ProtocolName:     "mkcp",
		ProtocolSettings: &Config{},
	}, func(conn internet.Connection) {
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
	})
	common.Must(err)
	defer listerner.Close()

	port := net.Port(listerner.Addr().(*net.UDPAddr).Port)

	var errg errgroup.Group
	for i := 0; i < 10; i++ {
		errg.Go(func() error {
			clientConn, err := DialKCP(context.Background(), net.UDPDestination(net.LocalHostIP, port), &internet.MemoryStreamConfig{
				ProtocolName:     "mkcp",
				ProtocolSettings: &Config{},
			})
			if err != nil {
				return err
			}
			defer clientConn.Close()

			clientSend := make([]byte, 1024*1024)
			rand.Read(clientSend)
			go clientConn.Write(clientSend)

			clientReceived := make([]byte, 1024*1024)
			common.Must2(io.ReadFull(clientConn, clientReceived))

			clientExpected := make([]byte, 1024*1024)
			for idx, b := range clientSend {
				clientExpected[idx] = b ^ 'c'
			}
			if r := cmp.Diff(clientReceived, clientExpected); r != "" {
				return errors.New(r)
			}
			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 60 && listerner.ActiveConnections() > 0; i++ {
		time.Sleep(500 * time.Millisecond)
	}
	if v := listerner.ActiveConnections(); v != 0 {
		t.Error("active connections: ", v)
	}
}
