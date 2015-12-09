package json_test

import (
	"encoding/json"
	"net"
	"testing"

	. "github.com/v2ray/v2ray-core/common/net/json"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestIPParsing(t *testing.T) {
	v2testing.Current(t)

	rawJson := "\"8.8.8.8\""
	host := &Host{}
	err := json.Unmarshal([]byte(rawJson), host)
	assert.Error(err).IsNil()
	assert.Bool(host.IsIP()).IsTrue()
	assert.Bool(host.IsDomain()).IsFalse()
	assert.Bool(host.IP().Equal(net.ParseIP("8.8.8.8"))).IsTrue()
}

func TestDomainParsing(t *testing.T) {
	v2testing.Current(t)

	rawJson := "\"v2ray.com\""
	host := &Host{}
	err := json.Unmarshal([]byte(rawJson), host)
	assert.Error(err).IsNil()
	assert.Bool(host.IsIP()).IsFalse()
	assert.Bool(host.IsDomain()).IsTrue()
	assert.StringLiteral(host.Domain()).Equals("v2ray.com")
}
