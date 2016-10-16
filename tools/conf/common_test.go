package conf_test

import (
	"encoding/json"
	"testing"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/tools/conf"
)

func TestIPParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := "\"8.8.8.8\""
	var address Address
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.Bytes([]byte(address.IP())).Equals([]byte{8, 8, 8, 8})
}

func TestDomainParsing(t *testing.T) {
	assert := assert.On(t)

	rawJson := "\"v2ray.com\""
	var address Address
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNil()
	assert.String(address.Domain()).Equals("v2ray.com")
}

func TestInvalidAddressJson(t *testing.T) {
	assert := assert.On(t)

	rawJson := "1234"
	var address Address
	err := json.Unmarshal([]byte(rawJson), &address)
	assert.Error(err).IsNotNil()
}

func TestStringNetwork(t *testing.T) {
	assert := assert.On(t)

	var network Network
	err := json.Unmarshal([]byte(`"tcp"`), &network)
	assert.Error(err).IsNil()
	assert.Bool(network.Build() == v2net.Network_TCP).IsTrue()
}

func TestArrayNetworkList(t *testing.T) {
	assert := assert.On(t)

	var list NetworkList
	err := json.Unmarshal([]byte("[\"Tcp\"]"), &list)
	assert.Error(err).IsNil()

	nlist := list.Build()
	assert.Bool(nlist.HasNetwork(v2net.ParseNetwork("tcp"))).IsTrue()
	assert.Bool(nlist.HasNetwork(v2net.ParseNetwork("udp"))).IsFalse()
}

func TestStringNetworkList(t *testing.T) {
	assert := assert.On(t)

	var list NetworkList
	err := json.Unmarshal([]byte("\"TCP, ip\""), &list)
	assert.Error(err).IsNil()

	nlist := list.Build()
	assert.Bool(nlist.HasNetwork(v2net.ParseNetwork("tcp"))).IsTrue()
	assert.Bool(nlist.HasNetwork(v2net.ParseNetwork("udp"))).IsFalse()
}

func TestInvalidNetworkJson(t *testing.T) {
	assert := assert.On(t)

	var list NetworkList
	err := json.Unmarshal([]byte("0"), &list)
	assert.Error(err).IsNotNil()
}
