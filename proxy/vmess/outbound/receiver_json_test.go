// +build json

package outbound_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
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
        "level": 255
      }
    ]
  }`

	receiver := new(Receiver)
	err := json.Unmarshal([]byte(rawJson), &receiver)
	assert.Error(err).IsNil()
	assert.String(receiver.Destination).Equals("tcp:127.0.0.1:80")
	assert.Int(len(receiver.Accounts)).Equals(1)
	assert.String(receiver.Accounts[0].ID).Equals("e641f5ad-9397-41e3-bf1a-e8740dfed019")
}
