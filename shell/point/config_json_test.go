// +build json

package point_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	_ "v2ray.com/core/app/router/rules"
	. "v2ray.com/core/shell/point"

	"v2ray.com/core/testing/assert"
)

func TestClientSampleConfig(t *testing.T) {
	assert := assert.On(t)

	GOPATH := os.Getenv("GOPATH")
	baseDir := filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core", "tools", "release", "config")

	pointConfig, err := LoadConfig(filepath.Join(baseDir, "vpoint_socks_vmess.json"))
	assert.Error(err).IsNil()

	assert.Pointer(pointConfig.InboundConfig).IsNotNil()
	assert.Port(pointConfig.InboundConfig.Port).IsValid()
	assert.Pointer(pointConfig.OutboundConfig).IsNotNil()

	assert.String(pointConfig.InboundConfig.Protocol).Equals("socks")
	assert.Pointer(pointConfig.InboundConfig.Settings).IsNotNil()

	assert.String(pointConfig.OutboundConfig.Protocol).Equals("vmess")
	assert.Pointer(pointConfig.OutboundConfig.Settings).IsNotNil()
}

func TestServerSampleConfig(t *testing.T) {
	assert := assert.On(t)

	GOPATH := os.Getenv("GOPATH")
	baseDir := filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core", "tools", "release", "config")

	pointConfig, err := LoadConfig(filepath.Join(baseDir, "vpoint_vmess_freedom.json"))
	assert.Error(err).IsNil()

	assert.Pointer(pointConfig.InboundConfig).IsNotNil()
	assert.Port(pointConfig.InboundConfig.Port).IsValid()
	assert.Pointer(pointConfig.OutboundConfig).IsNotNil()

	assert.String(pointConfig.InboundConfig.Protocol).Equals("vmess")
	assert.Pointer(pointConfig.InboundConfig.Settings).IsNotNil()

	assert.String(pointConfig.OutboundConfig.Protocol).Equals("freedom")
	assert.Pointer(pointConfig.OutboundConfig.Settings).IsNotNil()
}

func TestDefaultValueOfRandomAllocation(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "protocol": "vmess",
    "port": 1,
    "settings": {},
    "allocate": {
      "strategy": "random"
    }
  }`

	inboundDetourConfig := new(InboundDetourConfig)
	err := json.Unmarshal([]byte(rawJson), inboundDetourConfig)
	assert.Error(err).IsNil()
	assert.String(inboundDetourConfig.Allocation.Strategy).Equals(AllocationStrategyRandom)
	assert.Int(inboundDetourConfig.Allocation.Concurrency).Equals(3)
	assert.Int(inboundDetourConfig.Allocation.Refresh).Equals(5)
}
