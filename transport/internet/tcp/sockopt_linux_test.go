// +build linux

package tcp_test

import (
	"context"
	"strings"
	"testing"

	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	. "v2ray.com/core/transport/internet/tcp"
)

func TestGetOriginalDestination(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	conn, err := Dial(context.Background(), dest)
	assert.Error(err).IsNil()
	defer conn.Close()

	originalDest, err := GetOriginalDestination(conn)
	assert.Bool(dest == originalDest || strings.Contains(err.Error(), "failed to call getsockopt"))
}
