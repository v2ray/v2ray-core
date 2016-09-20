// +build json

package dns_test

import (
	"encoding/json"
	"testing"

	. "v2ray.com/core/app/dns"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
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
	dest := config.NameServers[0].AsDestination()
	assert.Destination(dest).IsUDP()
	assert.Address(dest.Address).Equals(v2net.IPAddress([]byte{8, 8, 8, 8}))
	assert.Port(dest.Port).Equals(v2net.Port(53))
}
