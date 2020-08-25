package core_test

import (
	"context"
	"crypto/rand"
	"io"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/testing/servers/udp"
)

func xor(b []byte) []byte {
	r := make([]byte, len(b))
	for i, v := range b {
		r[i] = v ^ 'c'
	}
	return r
}

func xor2(b []byte) []byte {
	r := make([]byte, len(b))
	for i, v := range b {
		r[i] = v ^ 'd'
	}
	return r
}

func TestV2RayDial(t *testing.T) {
	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	common.Must(err)
	defer tcpServer.Close()

	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	cfgBytes, err := proto.Marshal(config)
	common.Must(err)

	server, err := core.StartInstance("protobuf", cfgBytes)
	common.Must(err)
	defer server.Close()

	conn, err := core.Dial(context.Background(), server, dest)
	common.Must(err)
	defer conn.Close()

	const size = 10240 * 1024
	payload := make([]byte, size)
	common.Must2(rand.Read(payload))

	if _, err := conn.Write(payload); err != nil {
		t.Fatal(err)
	}

	receive := make([]byte, size)
	if _, err := io.ReadFull(conn, receive); err != nil {
		t.Fatal("failed to read all response: ", err)
	}

	if r := cmp.Diff(xor(receive), payload); r != "" {
		t.Error(r)
	}
}

func TestV2RayDialUDPConn(t *testing.T) {
	udpServer := udp.Server{
		MsgProcessor: xor,
	}
	dest, err := udpServer.Start()
	common.Must(err)
	defer udpServer.Close()

	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	cfgBytes, err := proto.Marshal(config)
	common.Must(err)

	server, err := core.StartInstance("protobuf", cfgBytes)
	common.Must(err)
	defer server.Close()

	conn, err := core.Dial(context.Background(), server, dest)
	common.Must(err)
	defer conn.Close()

	const size = 1024
	payload := make([]byte, size)
	common.Must2(rand.Read(payload))

	for i := 0; i < 2; i++ {
		if _, err := conn.Write(payload); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(time.Millisecond * 500)

	receive := make([]byte, size*2)
	for i := 0; i < 2; i++ {
		n, err := conn.Read(receive)
		if err != nil {
			t.Fatal("expect no error, but got ", err)
		}
		if n != size {
			t.Fatal("expect read size ", size, " but got ", n)
		}

		if r := cmp.Diff(xor(receive[:n]), payload); r != "" {
			t.Fatal(r)
		}
	}
}

func TestV2RayDialUDP(t *testing.T) {
	udpServer1 := udp.Server{
		MsgProcessor: xor,
	}
	dest1, err := udpServer1.Start()
	common.Must(err)
	defer udpServer1.Close()

	udpServer2 := udp.Server{
		MsgProcessor: xor2,
	}
	dest2, err := udpServer2.Start()
	common.Must(err)
	defer udpServer2.Close()

	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	cfgBytes, err := proto.Marshal(config)
	common.Must(err)

	server, err := core.StartInstance("protobuf", cfgBytes)
	common.Must(err)
	defer server.Close()

	conn, err := core.DialUDP(context.Background(), server)
	common.Must(err)
	defer conn.Close()

	const size = 1024
	{
		payload := make([]byte, size)
		common.Must2(rand.Read(payload))

		if _, err := conn.WriteTo(payload, &net.UDPAddr{
			IP:   dest1.Address.IP(),
			Port: int(dest1.Port),
		}); err != nil {
			t.Fatal(err)
		}

		receive := make([]byte, size)
		if _, _, err := conn.ReadFrom(receive); err != nil {
			t.Fatal(err)
		}

		if r := cmp.Diff(xor(receive), payload); r != "" {
			t.Error(r)
		}
	}

	{
		payload := make([]byte, size)
		common.Must2(rand.Read(payload))

		if _, err := conn.WriteTo(payload, &net.UDPAddr{
			IP:   dest2.Address.IP(),
			Port: int(dest2.Port),
		}); err != nil {
			t.Fatal(err)
		}

		receive := make([]byte, size)
		if _, _, err := conn.ReadFrom(receive); err != nil {
			t.Fatal(err)
		}

		if r := cmp.Diff(xor2(receive), payload); r != "" {
			t.Error(r)
		}
	}
}
