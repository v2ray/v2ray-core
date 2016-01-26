// +build json

package http_test

import (
	"encoding/json"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	. "github.com/v2ray/v2ray-core/proxy/http"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestOwnHosts(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "ownHosts": [
      "127.0.0.1",
      "google.com"
    ]
  }`

	config := new(Config)
	err := json.Unmarshal([]byte(rawJson), config)
	assert.Error(err).IsNil()
	assert.Bool(config.IsOwnHost(v2net.IPAddress([]byte{127, 0, 0, 1}))).IsTrue()
	assert.Bool(config.IsOwnHost(v2net.DomainAddress("google.com"))).IsTrue()
	assert.Bool(config.IsOwnHost(v2net.DomainAddress("local.v2ray.com"))).IsTrue()
	assert.Bool(config.IsOwnHost(v2net.DomainAddress("v2ray.com"))).IsFalse()
}
