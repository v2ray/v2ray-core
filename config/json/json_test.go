package json_test

import (
	"path/filepath"
	"testing"

	"github.com/v2ray/v2ray-core/config"
	"github.com/v2ray/v2ray-core/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/freedom/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/socks/config/json"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/config/json"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestClientSampleConfig(t *testing.T) {
	assert := unit.Assert(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	pointConfig, err := json.LoadConfig(filepath.Join(baseDir, "vpoint_socks_vmess.json"))
	assert.Error(err).IsNil()

	assert.Uint16(pointConfig.Port()).Positive()
	assert.Pointer(pointConfig.InboundConfig()).IsNotNil()
	assert.Pointer(pointConfig.OutboundConfig()).IsNotNil()

	assert.String(pointConfig.InboundConfig().Protocol()).Equals("socks")
	assert.Pointer(pointConfig.InboundConfig().Settings(config.TypeInbound)).IsNotNil()

	assert.String(pointConfig.OutboundConfig().Protocol()).Equals("vmess")
	assert.Pointer(pointConfig.OutboundConfig().Settings(config.TypeOutbound)).IsNotNil()
}

func TestServerSampleConfig(t *testing.T) {
	assert := unit.Assert(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	pointConfig, err := json.LoadConfig(filepath.Join(baseDir, "vpoint_vmess_freedom.json"))
	assert.Error(err).IsNil()

	assert.Uint16(pointConfig.Port()).Positive()
	assert.Pointer(pointConfig.InboundConfig()).IsNotNil()
	assert.Pointer(pointConfig.OutboundConfig()).IsNotNil()

	assert.String(pointConfig.InboundConfig().Protocol()).Equals("vmess")
	assert.Pointer(pointConfig.InboundConfig().Settings(config.TypeInbound)).IsNotNil()

	assert.String(pointConfig.OutboundConfig().Protocol()).Equals("freedom")
	assert.Pointer(pointConfig.OutboundConfig().Settings(config.TypeOutbound)).IsNotNil()
}
