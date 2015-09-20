package json

import (
	"path/filepath"
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestClientSampleConfig(t *testing.T) {
	assert := unit.Assert(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	config, err := LoadConfig(filepath.Join(baseDir, "vpoint_socks_vmess.json"))
	assert.Error(err).IsNil()

	assert.Uint16(config.Port()).Positive()
	assert.Pointer(config.InboundConfig()).IsNotNil()
	assert.Pointer(config.OutboundConfig()).IsNotNil()

	assert.String(config.InboundConfig().Protocol()).Equals("socks")
	assert.Int(len(config.InboundConfig().Content())).GreaterThan(0)

	assert.String(config.OutboundConfig().Protocol()).Equals("vmess")
	assert.Int(len(config.OutboundConfig().Content())).GreaterThan(0)
}

func TestServerSampleConfig(t *testing.T) {
	assert := unit.Assert(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	config, err := LoadConfig(filepath.Join(baseDir, "vpoint_vmess_freedom.json"))
	assert.Error(err).IsNil()

	assert.Uint16(config.Port()).Positive()
	assert.Pointer(config.InboundConfig()).IsNotNil()
	assert.Pointer(config.OutboundConfig()).IsNotNil()

	assert.String(config.InboundConfig().Protocol()).Equals("vmess")
	assert.Int(len(config.InboundConfig().Content())).GreaterThan(0)

	assert.String(config.OutboundConfig().Protocol()).Equals("freedom")
  assert.Int(len(config.OutboundConfig().Content())).Equals(0)
}
