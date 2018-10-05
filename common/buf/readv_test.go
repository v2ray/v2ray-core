// +build !wasm

package buf_test

import (
	"crypto/rand"
	"net"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/buf"
	"v2ray.com/core/common/compare"
	"v2ray.com/core/testing/servers/tcp"
)

func TestReadvReader(t *testing.T) {
	tcpServer := &tcp.Server{
		MsgProcessor: func(b []byte) []byte {
			return b
		},
	}
	dest, err := tcpServer.Start()
	common.Must(err)
	defer tcpServer.Close() // nolint: errcheck

	conn, err := net.Dial("tcp", dest.NetAddr())
	common.Must(err)
	defer conn.Close() // nolint: errcheck

	const size = 8192
	data := make([]byte, 8192)
	common.Must2(rand.Read(data))

	go func() {
		writer := NewWriter(conn)
		var mb MultiBuffer
		common.Must2(mb.Write(data))

		if err := writer.WriteMultiBuffer(mb); err != nil {
			t.Fatal("failed to write data: ", err)
		}
	}()

	rawConn, err := conn.(*net.TCPConn).SyscallConn()
	common.Must(err)

	reader := NewReadVReader(conn, rawConn)
	var rmb MultiBuffer
	for {
		mb, err := reader.ReadMultiBuffer()
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}
		rmb.AppendMulti(mb)
		if rmb.Len() == size {
			break
		}
	}

	rdata := make([]byte, size)
	common.Must2(rmb.Read(rdata))

	if err := compare.BytesEqualWithDetail(data, rdata); err != nil {
		t.Fatal(err)
	}
}
