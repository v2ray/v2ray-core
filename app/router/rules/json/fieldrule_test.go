package json

import (
	"encoding/json"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestStringListParsingList(t *testing.T) {
	v2testing.Current(t)

	rawJson := `["a", "b", "c", "d"]`
	var strList StringList
	err := json.Unmarshal([]byte(rawJson), &strList)
	assert.Error(err).IsNil()
	assert.Int(strList.Len()).Equals(4)
}

func TestStringListParsingString(t *testing.T) {
	v2testing.Current(t)

	rawJson := `"abcd"`
	var strList StringList
	err := json.Unmarshal([]byte(rawJson), &strList)
	assert.Error(err).IsNil()
	assert.Int(strList.Len()).Equals(1)
}

func TestDomainMatching(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "type": "field",
    "domain": ["google.com", "regexp:v2ray.com$"],
    "tag": "test"
  }`
	rule := parseRule([]byte(rawJson))
	dest := v2net.TCPDestination(v2net.DomainAddress("www.v2ray.com"), 80)
	assert.Bool(rule.Apply(dest)).IsTrue()
}

func TestPortMatching(t *testing.T) {
	v2testing.Current(t)

	rule := &FieldRule{
		Port: &v2nettesting.PortRange{
			FromValue: 0,
			ToValue:   100,
		},
	}
	dest := v2net.TCPDestination(v2net.DomainAddress("www.v2ray.com"), 80)
	assert.Bool(rule.Apply(dest)).IsTrue()
}

func TestIPMatching(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "type": "field",
    "ip": "10.0.0.0/8",
    "tag": "test"
  }`
	rule := parseRule([]byte(rawJson))
	dest := v2net.TCPDestination(v2net.IPAddress([]byte{10, 0, 0, 1}), 80)
	assert.Bool(rule.Apply(dest)).IsTrue()
}

func TestIPListMatching(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "type": "field",
    "ip": ["10.0.0.0/8", "192.168.0.0/16"],
    "tag": "test"
  }`
	rule := parseRule([]byte(rawJson))
	dest := v2net.TCPDestination(v2net.IPAddress([]byte{192, 168, 1, 1}), 80)
	assert.Bool(rule.Apply(dest)).IsTrue()
}

func TestPortNotMatching(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "type": "field",
    "port": "80-100",
    "tag": "test"
  }`
	rule := parseRule([]byte(rawJson))
	dest := v2net.TCPDestination(v2net.IPAddress([]byte{10, 0, 0, 1}), 79)
	assert.Bool(rule.Apply(dest)).IsFalse()
}

func TestDomainNotMatching(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "type": "field",
    "domain": ["google.com", "v2ray.com"],
    "tag": "test"
  }`
	rule := parseRule([]byte(rawJson))
	dest := v2net.TCPDestination(v2net.IPAddress([]byte{10, 0, 0, 1}), 80)
	assert.Bool(rule.Apply(dest)).IsFalse()
}

func TestDomainNotMatchingDomain(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "type": "field",
    "domain": ["google.com", "v2ray.com"],
    "tag": "test"
  }`
	rule := parseRule([]byte(rawJson))
	dest := v2net.TCPDestination(v2net.DomainAddress("baidu.com"), 80)
	assert.Bool(rule.Apply(dest)).IsFalse()

	dest = v2net.TCPDestination(v2net.DomainAddress("www.google.com"), 80)
	assert.Bool(rule.Apply(dest)).IsTrue()
}
