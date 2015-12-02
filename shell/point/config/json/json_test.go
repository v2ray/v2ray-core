package json_test

import (
	"path/filepath"
	"testing"

	_ "github.com/v2ray/v2ray-core/proxy/dokodemo/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/freedom/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/socks/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/config/json"
	"github.com/v2ray/v2ray-core/shell/point/config/json"

	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestClientSampleConfig(t *testing.T) {
	v2testing.Current(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	pointConfig, err := json.LoadConfig(filepath.Join(baseDir, "vpoint_socks_vmess.json"))
	assert.Error(err).IsNil()

	assert.Uint16(pointConfig.Port()).Positive()
	assert.Pointer(pointConfig.InboundConfig()).IsNotNil()
	assert.Pointer(pointConfig.OutboundConfig()).IsNotNil()

	assert.StringLiteral(pointConfig.InboundConfig().Protocol()).Equals("socks")
	assert.Pointer(pointConfig.InboundConfig().Settings()).IsNotNil()

	assert.StringLiteral(pointConfig.OutboundConfig().Protocol()).Equals("vmess")
	assert.Pointer(pointConfig.OutboundConfig().Settings()).IsNotNil()
}

func TestServerSampleConfig(t *testing.T) {
	v2testing.Current(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	pointConfig, err := json.LoadConfig(filepath.Join(baseDir, "vpoint_vmess_freedom.json"))
	assert.Error(err).IsNil()

	assert.Uint16(pointConfig.Port()).Positive()
	assert.Pointer(pointConfig.InboundConfig()).IsNotNil()
	assert.Pointer(pointConfig.OutboundConfig()).IsNotNil()

	assert.StringLiteral(pointConfig.InboundConfig().Protocol()).Equals("vmess")
	assert.Pointer(pointConfig.InboundConfig().Settings()).IsNotNil()

	assert.StringLiteral(pointConfig.OutboundConfig().Protocol()).Equals("freedom")
	assert.Pointer(pointConfig.OutboundConfig().Settings()).IsNotNil()
}

func TestDetourConfig(t *testing.T) {
	v2testing.Current(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	pointConfig, err := json.LoadConfig(filepath.Join(baseDir, "vpoint_dns_detour.json"))
	assert.Error(err).IsNil()

	detours := pointConfig.InboundDetours()
	assert.Int(len(detours)).Equals(1)

	detour := detours[0]
	assert.StringLiteral(detour.Protocol()).Equals("dokodemo-door")
	assert.Uint16(detour.PortRange().From().Value()).Equals(uint16(28394))
	assert.Uint16(detour.PortRange().To().Value()).Equals(uint16(28394))
	assert.Pointer(detour.Settings()).IsNotNil()
}
