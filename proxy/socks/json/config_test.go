package json

import (
	"encoding/json"
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/proxy/common/config"
	jsonconfig "github.com/v2ray/v2ray-core/proxy/common/config/json"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestAccountMapParsing(t *testing.T) {
	v2testing.Current(t)

	var accountMap SocksAccountMap
	err := json.Unmarshal([]byte("[{\"user\": \"a\", \"pass\":\"b\"}, {\"user\": \"c\", \"pass\":\"d\"}]"), &accountMap)
	assert.Error(err).IsNil()

	assert.Bool(accountMap.HasAccount("a", "b")).IsTrue()
	assert.Bool(accountMap.HasAccount("a", "c")).IsFalse()
	assert.Bool(accountMap.HasAccount("c", "d")).IsTrue()
	assert.Bool(accountMap.HasAccount("e", "d")).IsFalse()
}

func TestDefaultIPAddress(t *testing.T) {
	v2testing.Current(t)

	socksConfig := jsonconfig.CreateConfig("socks", config.TypeInbound).(*SocksConfig)
	assert.String(socksConfig.IP()).Equals("127.0.0.1")
}

func TestIPAddressParsing(t *testing.T) {
	v2testing.Current(t)

	var ipAddress IPAddress
	err := json.Unmarshal([]byte("\"1.2.3.4\""), &ipAddress)
	assert.Error(err).IsNil()
	assert.String(net.IP(ipAddress)).Equals("1.2.3.4")
}

func TestNoAuthConfig(t *testing.T) {
	v2testing.Current(t)

	var config SocksConfig
	err := json.Unmarshal([]byte("{\"auth\":\"noauth\", \"ip\":\"8.8.8.8\"}"), &config)
	assert.Error(err).IsNil()
	assert.Bool(config.IsNoAuth()).IsTrue()
	assert.Bool(config.IsPassword()).IsFalse()
	assert.String(config.IP()).Equals("8.8.8.8")
	assert.Bool(config.UDPEnabled()).IsFalse()
}

func TestUserPassConfig(t *testing.T) {
	v2testing.Current(t)

	var config SocksConfig
	err := json.Unmarshal([]byte("{\"auth\":\"password\", \"accounts\":[{\"user\":\"x\", \"pass\":\"y\"}], \"udp\":true}"), &config)
	assert.Error(err).IsNil()
	assert.Bool(config.IsNoAuth()).IsFalse()
	assert.Bool(config.IsPassword()).IsTrue()
	assert.Bool(config.HasAccount("x", "y")).IsTrue()
	assert.Bool(config.UDPEnabled()).IsTrue()
}
