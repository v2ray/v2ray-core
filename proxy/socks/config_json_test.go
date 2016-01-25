// +build json

package socks_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/proxy/internal/config"
	"github.com/v2ray/v2ray-core/proxy/socks"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestDefaultIPAddress(t *testing.T) {
	v2testing.Current(t)

	socksConfig, err := config.CreateInboundConfig("socks", []byte(`{
    "auth": "noauth"
  }`))
	assert.Error(err).IsNil()
	assert.String(socksConfig.(*socks.Config).Address).Equals("127.0.0.1")
}
