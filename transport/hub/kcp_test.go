package hub_test

import "testing"

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/transport"
	"github.com/v2ray/v2ray-core/transport/hub"
	"github.com/v2ray/v2ray-core/transport/hub/kcpv"
)

func Test_Pair(t *testing.T) {
	assert := assert.On(t)
	transport.KcpConfig = &kcpv.Config{}
	transport.KcpConfig.Mode = "fast2"
	transport.KcpConfig.Key = "key"
	transport.KcpConfig.AdvancedConfigs = kcpv.DefaultAdvancedConfigs
	lst, _ := hub.ListenKCP(v2net.ParseAddress("127.0.0.1"), 17777)
	go func() {
		connx, err2 := lst.Accept()
		assert.Error(err2).IsNil()
		connx.Close()
	}()
	conn, _ := hub.DialKCP(v2net.TCPDestination(v2net.ParseAddress("127.0.0.1"), 17777))
	conn.LocalAddr()
	conn.RemoteAddr()
	conn.ApplyConf()
	conn.Write([]byte("x"))
	conn.Close()
}
