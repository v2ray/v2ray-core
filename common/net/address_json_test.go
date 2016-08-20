// +build json

package net_test

import (
	"encoding/json"
	"net"
	"testing"

	. "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestIPParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := "\"8.8.8.8\""
	var address AddressJson
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.Bool(address.Address.Family().Either(AddressFamilyIPv4)).IsTrue()
	assert.Bool(address.Address.Family().Either(AddressFamilyDomain)).IsFalse()
	assert.Bool(address.Address.IP().Equal(net.ParseIP("8.8.8.8"))).IsTrue()
}

func TestDomainParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := "\"v2ray.com\""
	var address AddressJson
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.Bool(address.Address.Family().Either(AddressFamilyIPv4)).IsFalse()
	assert.Bool(address.Address.Family().Either(AddressFamilyDomain)).IsTrue()
	assert.String(address.Address.Domain()).Equals("v2ray.com")
}

func TestInvalidAddressJson(t *testing.T) {
	assert := assert.On(t)

	rawJson := "1234"
	var address AddressJson
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNotNil()
}
