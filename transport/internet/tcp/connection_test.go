package tcp_test

import (
	"net"
	"testing"

	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/tcp"
)

func TestRawConnection(t *testing.T) {
	assert := assert.On(t)

	rawConn := RawConnection{net.TCPConn{}}
	assert.Bool(rawConn.Reusable()).IsFalse()

	rawConn.SetReusable(true)
	assert.Bool(rawConn.Reusable()).IsFalse()
}
