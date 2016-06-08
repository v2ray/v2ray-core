// +build json

package dns_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/app/dns"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestConfigParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := `{
    "servers": ["8.8.8.8"]
  }`

	config := new(Config)
	err := json.Unmarshal([]byte(rawJson), config)
	assert.Error(err).IsNil()
	assert.Int(len(config.NameServers)).Equals(1)
	assert.Destination(config.NameServers[0]).IsUDP()
	assert.Address(config.NameServers[0].Address()).Equals(v2net.IPAddress([]byte{8, 8, 8, 8}))
	assert.Port(config.NameServers[0].Port()).Equals(v2net.Port(53))
}
