// +build linux

package tcp_test

import (
	"context"
	"strings"
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/transport/internet"
	. "v2ray.com/core/transport/internet/tcp"
	. "v2ray.com/ext/assert"
)

func TestGetOriginalDestination(t *testing.T) {
	assert := With(t)

	tcpServer := tcp.Server{}
	dest, err := tcpServer.Start()
	common.Must(err)
	defer tcpServer.Close()

	config, err := internet.ToMemoryStreamConfig(nil)
	common.Must(err)
	conn, err := Dial(context.Background(), dest, config)
	common.Must(err)
	defer conn.Close()

	originalDest, err := GetOriginalDestination(conn)
	assert(dest == originalDest || strings.Contains(err.Error(), "failed to call getsockopt"), IsTrue)
}
