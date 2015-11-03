package json

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/proxy/common/config"
	jsonconfig "github.com/v2ray/v2ray-core/proxy/common/config/json"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestAccountMapParsing(t *testing.T) {
	assert := unit.Assert(t)

	var accountMap SocksAccountMap
	err := json.Unmarshal([]byte("[{\"user\": \"a\", \"pass\":\"b\"}, {\"user\": \"c\", \"pass\":\"d\"}]"), &accountMap)
	assert.Error(err).IsNil()

	value, found := accountMap["a"]
	assert.Bool(found).IsTrue()
	assert.String(value).Equals("b")

	value, found = accountMap["c"]
	assert.Bool(found).IsTrue()
	assert.String(value).Equals("d")
}

func TestDefaultIPAddress(t *testing.T) {
	assert := unit.Assert(t)

	socksConfig := jsonconfig.CreateConfig("socks", config.TypeInbound).(*SocksConfig)
	assert.String(socksConfig.IP().String()).Equals("127.0.0.1")
}

func TestIPAddressParsing(t *testing.T) {
	assert := unit.Assert(t)

	var ipAddress IPAddress
	err := json.Unmarshal([]byte("\"1.2.3.4\""), &ipAddress)
	assert.Error(err).IsNil()
	assert.String(net.IP(ipAddress).String()).Equals("1.2.3.4")
}
