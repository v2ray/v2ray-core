package internet_test

import (
	"testing"

	"v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	. "v2ray.com/core/transport/internet"
)

func TestDialWithLocalAddr(t *testing.T) {
	assert := assert.On(t)

	server := &tcp.Server{}
	dest, err := server.Start()
	assert.Error(err).IsNil()
	defer server.Close()

	conn, err := DialSystem(net.LocalHostIP, net.TCPDestination(net.LocalHostIP, dest.Port))
	assert.Error(err).IsNil()
	assert.String(conn.RemoteAddr().String()).Equals("127.0.0.1:" + dest.Port.String())
	conn.Close()
}
