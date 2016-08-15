// +build linux

package udp_test

import (
	"runtime"
	"syscall"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/transport/internet"
	. "github.com/v2ray/v2ray-core/transport/internet/udp"
)

func TestHubSocksOption(t *testing.T) {
	assert := assert.On(t)

	hub, err := ListenUDP(v2net.LocalHostIP, v2net.Port(0), ListenOption{
		Callback:            func(*alloc.Buffer, *proxy.SessionInfo) {},
		ReceiveOriginalDest: true,
	})
	assert.Error(err).IsNil()
	conn := hub.Connection()

	sysfd, ok := conn.(internet.SysFd)
	assert.Bool(ok).IsTrue()

	fd, err := sysfd.SysFd()
	assert.Error(err).IsNil()

	v, err := syscall.GetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT)
	assert.Error(err).IsNil()
	assert.Int(v).Equals(1)

	v, err = syscall.GetsockoptInt(fd, syscall.SOL_IP, syscall.IP_RECVORIGDSTADDR)
	assert.Error(err).IsNil()
	assert.Int(v).Equals(1)
}
