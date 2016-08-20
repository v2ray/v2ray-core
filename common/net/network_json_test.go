// +build json

package net_test

import (
	"encoding/json"
	"testing"

	. "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestArrayNetworkList(t *testing.T) {
	assert := assert.On(t)

	var list NetworkList
	err := json.Unmarshal([]byte("[\"Tcp\"]"), &list)
	assert.Error(err).IsNil()
	assert.Bool(list.HasNetwork(Network("tcp"))).IsTrue()
	assert.Bool(list.HasNetwork(Network("udp"))).IsFalse()
}

func TestStringNetworkList(t *testing.T) {
	assert := assert.On(t)

	var list NetworkList
	err := json.Unmarshal([]byte("\"TCP, ip\""), &list)
	assert.Error(err).IsNil()
	assert.Bool(list.HasNetwork(Network("tcp"))).IsTrue()
	assert.Bool(list.HasNetwork(Network("udp"))).IsFalse()
}

func TestInvalidNetworkJson(t *testing.T) {
	assert := assert.On(t)

	var list NetworkList
	err := json.Unmarshal([]byte("0"), &list)
	assert.Error(err).IsNotNil()
}
