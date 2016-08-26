// +build json

package net_test

import (
	"encoding/json"
	"testing"

	. "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestIPParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := "\"8.8.8.8\""
	var address AddressPB
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.Bytes(address.GetIp()).Equals([]byte{8, 8, 8, 8})
}

func TestDomainParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := "\"v2ray.com\""
	var address AddressPB
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.String(address.GetDomain()).Equals("v2ray.com")
}

func TestInvalidAddressJson(t *testing.T) {
	assert := assert.On(t)

	rawJson := "1234"
	var address AddressPB
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNotNil()
}
