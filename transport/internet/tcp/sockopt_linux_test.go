// +build linux

package tcp_test

import (
	"context"
	"testing"

	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/tcp"
)

func TestGetOriginalDestination(t *testing.T) {
	assert := assert.On(t)

	tcpServer := tcp.Server{}
	dest, err := tcpServer.Start()
	assert.Error(err).IsNil()
	defer tcpServer.Close()

	ctx := internet.ContextWithTransportSettings(context.Background, &Config{})
	conn, err := Dial(ctx, dest)
	assert.Error(err).IsNil()
	defer conn.Close()

	_, err = GetOriginalDestination(conn)
	assert.String(err.Error()).Contains("failed to call getsockopt")
}
