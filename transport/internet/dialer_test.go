package internet_test

import (
	"context"
	"testing"

	"v2ray.com/core/common/net"
	"v2ray.com/core/testing/servers/tcp"
	. "v2ray.com/core/transport/internet"
	. "v2ray.com/ext/assert"
)

func TestDialWithLocalAddr(t *testing.T) {
	assert := With(t)

	server := &tcp.Server{}
	dest, err := server.Start()
	assert(err, IsNil)
	defer server.Close()

	conn, err := DialSystem(context.Background(), net.TCPDestination(net.LocalHostIP, dest.Port))
	assert(err, IsNil)
	assert(conn.RemoteAddr().String(), Equals, "127.0.0.1:"+dest.Port.String())
	conn.Close()
}
