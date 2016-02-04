// +build json

package net_test

import (
	"encoding/json"
	"net"
	"testing"

	. "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestIPParsing(t *testing.T) {
	v2testing.Current(t)

	rawJson := "\"8.8.8.8\""
	var address AddressJson
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.Bool(address.Address.IsIPv4()).IsTrue()
	assert.Bool(address.Address.IsDomain()).IsFalse()
	assert.Bool(address.Address.IP().Equal(net.ParseIP("8.8.8.8"))).IsTrue()
}

func TestDomainParsing(t *testing.T) {
	v2testing.Current(t)

	rawJson := "\"v2ray.com\""
	var address AddressJson
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.Bool(address.Address.IsIPv4()).IsFalse()
	assert.Bool(address.Address.IsDomain()).IsTrue()
	assert.StringLiteral(address.Address.Domain()).Equals("v2ray.com")
}

func TestInvalidJson(t *testing.T) {
	v2testing.Current(t)

	rawJson := "1234"
	var address AddressJson
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNotNil()
}
