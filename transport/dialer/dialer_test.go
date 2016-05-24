package dialer_test

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	. "github.com/v2ray/v2ray-core/transport/dialer"
)

func TestDialDomain(t *testing.T) {
	assert := assert.On(t)

	server := &tcp.Server{
		Port: v2nettesting.PickPort(),
	}
	dest, err := server.Start()
	assert.Error(err).IsNil()
	defer server.Close()

	conn, err := Dial(v2net.TCPDestination(v2net.DomainAddress("local.v2ray.com"), dest.Port()))
	assert.Error(err).IsNil()
	assert.String(conn.RemoteAddr().String()).Equals("127.0.0.1:" + dest.Port().String())
	conn.Close()
}
