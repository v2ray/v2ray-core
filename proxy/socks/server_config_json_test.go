// +build json

package socks_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/proxy/registry"
	"github.com/v2ray/v2ray-core/proxy/socks"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestDefaultIPAddress(t *testing.T) {
	assert := assert.On(t)

	socksConfig, err := registry.CreateInboundConfig("socks", []byte(`{
    "auth": "noauth"
  }`))
	assert.Error(err).IsNil()
	assert.Address(socksConfig.(*socks.Config).Address).EqualsString("127.0.0.1")
}
