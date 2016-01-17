package router_test

import (
	"net"
	"path/filepath"
	"testing"

	. "github.com/v2ray/v2ray-core/app/router"
	_ "github.com/v2ray/v2ray-core/app/router/rules"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/shell/point"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestRouter(t *testing.T) {
	v2testing.Current(t)

	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	pointConfig, err := point.LoadConfig(filepath.Join(baseDir, "vpoint_socks_vmess.json"))
	assert.Error(err).IsNil()

	router, err := CreateRouter(pointConfig.RouterConfig.Strategy, pointConfig.RouterConfig.Settings)
	assert.Error(err).IsNil()

	dest := v2net.TCPDestination(v2net.IPAddress(net.ParseIP("120.135.126.1")), 80)
	tag, err := router.TakeDetour(dest)
	assert.Error(err).IsNil()
	assert.StringLiteral(tag).Equals("direct")
}
