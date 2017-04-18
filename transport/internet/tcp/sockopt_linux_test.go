// +build linux

package tcp_test

import (
	"context"
	"testing"

	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
)

func TestGetOriginalDestination(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{
		MsgProcessor: xor,
	}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	conn, err := Dial(context.Background(), dest)
	assert.Error(err).IsNil()

	_, err := GetOriginalDestination(conn)
	assert.String(err.Error()).Contains("failed to call getsockopt")
}
