package json

import (
	"encoding/json"
	"testing"

	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestArrayNetworkList(t *testing.T) {
	v2testing.Current(t)

	var list NetworkList
	err := json.Unmarshal([]byte("[\"Tcp\"]"), &list)
	assert.Error(err).IsNil()
	assert.Bool(list.HasNetwork("tcp")).IsTrue()
	assert.Bool(list.HasNetwork("udp")).IsFalse()
}

func TestStringNetworkList(t *testing.T) {
	v2testing.Current(t)

	var list NetworkList
	err := json.Unmarshal([]byte("\"TCP, ip\""), &list)
	assert.Error(err).IsNil()
	assert.Bool(list.HasNetwork("tcp")).IsTrue()
	assert.Bool(list.HasNetwork("udp")).IsFalse()
}
