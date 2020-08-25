// +build !windows

package domainsocket_test

import (
	"context"
	"runtime"
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/domainsocket"
)

func TestListen(t *testing.T) {
	ctx := context.Background()
	streamSettings := &internet.MemoryStreamConfig{
		ProtocolName: "domainsocket",
		ProtocolSettings: &Config{
			Path: "/tmp/ts3",
		},
	}
	listener, err := Listen(ctx, nil, net.Port(0), streamSettings, func(conn internet.Connection) {
		defer conn.Close()

		b := buf.New()
		defer b.Release()
		common.Must2(b.ReadFrom(conn))
		b.WriteString("Response")

		common.Must2(conn.Write(b.Bytes()))
	})
	common.Must(err)
	defer listener.Close()

	conn, err := Dial(ctx, net.Destination{}, streamSettings)
	common.Must(err)
	defer conn.Close()

	common.Must2(conn.Write([]byte("Request")))

	b := buf.New()
	defer b.Release()
	common.Must2(b.ReadFrom(conn))

	if b.String() != "RequestResponse" {
		t.Error("expected response as 'RequestResponse' but got ", b.String())
	}
}

func TestListenAbstract(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	ctx := context.Background()
	streamSettings := &internet.MemoryStreamConfig{
		ProtocolName: "domainsocket",
		ProtocolSettings: &Config{
			Path:     "/tmp/ts3",
			Abstract: true,
		},
	}
	listener, err := Listen(ctx, nil, net.Port(0), streamSettings, func(conn internet.Connection) {
		defer conn.Close()

		b := buf.New()
		defer b.Release()
		common.Must2(b.ReadFrom(conn))
		b.WriteString("Response")

		common.Must2(conn.Write(b.Bytes()))
	})
	common.Must(err)
	defer listener.Close()

	conn, err := Dial(ctx, net.Destination{}, streamSettings)
	common.Must(err)
	defer conn.Close()

	common.Must2(conn.Write([]byte("Request")))

	b := buf.New()
	defer b.Release()
	common.Must2(b.ReadFrom(conn))

	if b.String() != "RequestResponse" {
		t.Error("expected response as 'RequestResponse' but got ", b.String())
	}
}
