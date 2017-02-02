package conf_test

import (
	"encoding/json"
	"testing"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/tools/conf"
)

func TestStringListUnmarshalError(t *testing.T) {
	assert := assert.On(t)

	rawJson := `1234`
	list := new(StringList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert.Error(err).IsNotNil()
}

func TestStringListLen(t *testing.T) {
	assert := assert.On(t)

	rawJson := `"a, b, c, d"`
	list := new(StringList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert.Error(err).IsNil()
	assert.Int(list.Len()).Equals(4)
}

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

func TestIntPort(t *testing.T) {
	assert := assert.On(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("1234"), &portRange)
	assert.Error(err).IsNil()

	assert.Uint32(portRange.From).Equals(1234)
	assert.Uint32(portRange.To).Equals(1234)
}

func TestOverRangeIntPort(t *testing.T) {
	assert := assert.On(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("70000"), &portRange)
	assert.Error(err).IsNotNil()

	err = json.Unmarshal([]byte("-1"), &portRange)
	assert.Error(err).IsNotNil()
}

func TestSingleStringPort(t *testing.T) {
	assert := assert.On(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234\""), &portRange)
	assert.Error(err).IsNil()

	assert.Uint32(portRange.From).Equals(1234)
	assert.Uint32(portRange.To).Equals(1234)
}

func TestStringPairPort(t *testing.T) {
	assert := assert.On(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"1234-5678\""), &portRange)
	assert.Error(err).IsNil()

	assert.Uint32(portRange.From).Equals(1234)
	assert.Uint32(portRange.To).Equals(5678)
}

func TestOverRangeStringPort(t *testing.T) {
	assert := assert.On(t)

	var portRange PortRange
	err := json.Unmarshal([]byte("\"65536\""), &portRange)
	assert.Error(err).IsNotNil()

	err = json.Unmarshal([]byte("\"70000-80000\""), &portRange)
	assert.Error(err).IsNotNil()

	err = json.Unmarshal([]byte("\"1-90000\""), &portRange)
	assert.Error(err).IsNotNil()

	err = json.Unmarshal([]byte("\"700-600\""), &portRange)
	assert.Error(err).IsNotNil()
}

func TestUserParsing(t *testing.T) {
	assert := assert.On(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{
    "id": "96edb838-6d68-42ef-a933-25f7ac3a9d09",
    "email": "love@v2ray.com",
    "level": 1,
    "alterId": 100
  }`), user)
	assert.Error(err).IsNil()

	nUser := user.Build()
	assert.Byte(byte(nUser.Level)).Equals(1)
	assert.String(nUser.Email).Equals("love@v2ray.com")
}

func TestInvalidUserJson(t *testing.T) {
	assert := assert.On(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{"email": 1234}`), user)
	assert.Error(err).IsNotNil()
}
