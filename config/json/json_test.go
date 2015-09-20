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

	assert.Uint16(config.PortValue).Positive()
	assert.Pointer(config.InboundConfigValue).IsNotNil()
	assert.Pointer(config.OutboundConfigValue).IsNotNil()

	assert.String(config.InboundConfigValue.ProtocolString).Equals("socks")
	assert.Int(len(config.InboundConfigValue.Content())).GreaterThan(0)

	assert.String(config.OutboundConfigValue.ProtocolString).Equals("vmess")
	assert.Int(len(config.OutboundConfigValue.Content())).GreaterThan(0)
}

func TestServerSampleConfig(t *testing.T) {
	assert := unit.Assert(t)

	// TODO: fix for Windows
	baseDir := "$GOPATH/src/github.com/v2ray/v2ray-core/release/config"

	config, err := LoadConfig(filepath.Join(baseDir, "vpoint_vmess_freedom.json"))
	assert.Error(err).IsNil()

	assert.Uint16(config.PortValue).Positive()
	assert.Pointer(config.InboundConfigValue).IsNotNil()
	assert.Pointer(config.OutboundConfigValue).IsNotNil()

	assert.String(config.InboundConfigValue.ProtocolString).Equals("vmess")
	assert.Int(len(config.InboundConfigValue.Content())).GreaterThan(0)

	assert.String(config.OutboundConfigValue.ProtocolString).Equals("freedom")
}
