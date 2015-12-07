package json_test

import (
	"encoding/json"
	"testing"

	jsonconfig "github.com/v2ray/v2ray-core/proxy/vmess/outbound/json"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestConfigTargetParsing(t *testing.T) {
	v2testing.Current(t)

	rawJson := `{
    "address": "127.0.0.1",
    "port": 80,
    "users": [
      {
        "id": "e641f5ad-9397-41e3-bf1a-e8740dfed019",
        "email": "love@v2ray.com",
        "level": 999
      }
    ]
  }`

	var target *jsonconfig.ConfigTarget
	err := json.Unmarshal([]byte(rawJson), &target)
	assert.Error(err).IsNil()
	assert.String(target.Address).Equals("127.0.0.1:80")
	assert.Int(len(target.Users)).Equals(1)
	assert.StringLiteral(target.Users[0].ID().String).Equals("e641f5ad-9397-41e3-bf1a-e8740dfed019")
}
