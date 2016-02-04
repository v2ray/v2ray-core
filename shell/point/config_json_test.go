// +build json

package point_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/v2ray/v2ray-core/app/router/rules"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	. "github.com/v2ray/v2ray-core/shell/point"

	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestClientSampleConfig(t *testing.T) {
	v2testing.Current(t)

	GOPATH := os.Getenv("GOPATH")
	baseDir := filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core", "release", "config")

	pointConfig, err := LoadConfig(filepath.Join(baseDir, "vpoint_socks_vmess.json"))
	assert.Error(err).IsNil()

	netassert.Port(pointConfig.Port).IsValid()
	assert.Pointer(pointConfig.InboundConfig).IsNotNil()
	assert.Pointer(pointConfig.OutboundConfig).IsNotNil()

	assert.StringLiteral(pointConfig.InboundConfig.Protocol).Equals("socks")
	assert.Pointer(pointConfig.InboundConfig.Settings).IsNotNil()

	assert.StringLiteral(pointConfig.OutboundConfig.Protocol).Equals("vmess")
	assert.Pointer(pointConfig.OutboundConfig.Settings).IsNotNil()
}

func TestServerSampleConfig(t *testing.T) {
	v2testing.Current(t)

	GOPATH := os.Getenv("GOPATH")
	baseDir := filepath.Join(GOPATH, "src", "github.com", "v2ray", "v2ray-core", "release", "config")

	pointConfig, err := LoadConfig(filepath.Join(baseDir, "vpoint_vmess_freedom.json"))
	assert.Error(err).IsNil()

	netassert.Port(pointConfig.Port).IsValid()
	assert.Pointer(pointConfig.InboundConfig).IsNotNil()
	assert.Pointer(pointConfig.OutboundConfig).IsNotNil()

	assert.StringLiteral(pointConfig.InboundConfig.Protocol).Equals("vmess")
	assert.Pointer(pointConfig.InboundConfig.Settings).IsNotNil()

	assert.StringLiteral(pointConfig.OutboundConfig.Protocol).Equals("freedom")
	assert.Pointer(pointConfig.OutboundConfig.Settings).IsNotNil()
}

func TestDefaultValueOfRandomAllocation(t *testing.T) {
	v2testing.Current(t)

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
	assert.StringLiteral(inboundDetourConfig.Allocation.Strategy).Equals(AllocationStrategyRandom)
	assert.Int(inboundDetourConfig.Allocation.Concurrency).Equals(3)
	assert.Int(inboundDetourConfig.Allocation.Refresh).Equals(5)
}
