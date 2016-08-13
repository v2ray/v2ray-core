package ws_test

import (
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/tcp"
)

func TestRawConnection(t *testing.T) {
	assert := assert.On(t)

	rawConn := RawConnection{net.TCPConn{}}
	assert.Bool(rawConn.Reusable()).IsFalse()

	rawConn.SetReusable(true)
	assert.Bool(rawConn.Reusable()).IsFalse()
}
