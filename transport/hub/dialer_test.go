package hub_test

import (
	"net"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	. "github.com/v2ray/v2ray-core/transport/hub"
)

func TestDialDomain(t *testing.T) {
	assert := assert.On(t)

	server := &tcp.Server{
		Port: v2nettesting.PickPort(),
	}
	dest, err := server.Start()
	assert.Error(err).IsNil()
	defer server.Close()

	conn, err := Dial(nil, v2net.TCPDestination(v2net.DomainAddress("local.v2ray.com"), dest.Port()))
	assert.Error(err).IsNil()
	assert.String(conn.RemoteAddr().String()).Equals("127.0.0.1:" + dest.Port().String())
	conn.Close()
}

func TestDialWithLocalAddr(t *testing.T) {
	assert := assert.On(t)

	server := &tcp.Server{
		Port: v2nettesting.PickPort(),
	}
	dest, err := server.Start()
	assert.Error(err).IsNil()
	defer server.Close()

	var localAddr net.IP
	addrs, err := net.InterfaceAddrs()
	assert.Error(err).IsNil()
	for _, addr := range addrs {
		str := addr.String()
		ip := net.ParseIP(str)
		if ip != nil && ip.To4() != nil {
			localAddr = ip.To4()
		}
	}
	assert.Pointer(localAddr).IsNotNil()

	conn, err := Dial(v2net.IPAddress(localAddr), v2net.TCPDestination(v2net.LocalHostIP, dest.Port()))
	assert.Error(err).IsNil()
	assert.String(conn.RemoteAddr().String()).Equals("127.0.0.1:" + dest.Port().String())
	conn.Close()
}
