// +build json

package socks_test

import (
	"testing"

	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/proxy/socks"
	"v2ray.com/core/testing/assert"
)

func TestDefaultIPAddress(t *testing.T) {
	assert := assert.On(t)

	socksConfig, err := registry.CreateInboundConfig("socks", []byte(`{
    "auth": "noauth"
  }`))
	assert.Error(err).IsNil()
	assert.Address(socksConfig.(*socks.Config).Address).EqualsString("127.0.0.1")
}
