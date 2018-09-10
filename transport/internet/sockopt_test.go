package internet_test

import (
	"context"
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/compare"
	"v2ray.com/core/testing/servers/tcp"
	. "v2ray.com/core/transport/internet"
)

func TestSockOptMark(t *testing.T) {
	tcpServer := tcp.Server{
		MsgProcessor: func(b []byte) []byte {
			return b
		},
	}
	dest, err := tcpServer.Start()
	common.Must(err)
	defer tcpServer.Close()

	ctx := context.Background()
	ctx = ContextWithStreamSettings(ctx, &MemoryStreamConfig{
		SocketSettings: &SocketConfig{
			Tfo: SocketConfig_Enable,
		},
	})
	dialer := DefaultSystemDialer{}
	conn, err := dialer.Dial(ctx, nil, dest)
	common.Must(err)
	defer conn.Close()

	_, err = conn.Write([]byte("abcd"))
	common.Must(err)

	b := buf.New()
	common.Must(b.Reset(buf.ReadFrom(conn)))
	if err := compare.BytesEqualWithDetail(b.Bytes(), []byte("abcd")); err != nil {
		t.Fatal(err)
	}
}
