package internet_test

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	. "github.com/v2ray/v2ray-core/transport/internet"
)

func TestDialDomain(t *testing.T) {
	assert := assert.On(t)

	server := &tcp.Server{}
	dest, err := server.Start()
	assert.Error(err).IsNil()
	defer server.Close()

	conn, err := DialToDest(nil, v2net.TCPDestination(v2net.DomainAddress("local.v2ray.com"), dest.Port()))
	assert.Error(err).IsNil()
	assert.String(conn.RemoteAddr().String()).Equals("127.0.0.1:" + dest.Port().String())
	conn.Close()
}

func TestDialWithLocalAddr(t *testing.T) {
	assert := assert.On(t)

	server := &tcp.Server{}
	dest, err := server.Start()
	assert.Error(err).IsNil()
	defer server.Close()

	conn, err := DialToDest(v2net.LocalHostIP, v2net.TCPDestination(v2net.LocalHostIP, dest.Port()))
	assert.Error(err).IsNil()
	assert.String(conn.RemoteAddr().String()).Equals("127.0.0.1:" + dest.Port().String())
	conn.Close()
}
