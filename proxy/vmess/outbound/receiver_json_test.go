// +build json

package outbound_test

import (
	"encoding/json"
	"testing"

	"github.com/v2ray/v2ray-core/common/protocol"
	. "github.com/v2ray/v2ray-core/proxy/vmess/outbound"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestConfigTargetParsing(t *testing.T) {
	assert := assert.On(t)

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
	assert.Destination(receiver.Destination).EqualsString("tcp:127.0.0.1:80")
	assert.Int(len(receiver.Accounts)).Equals(1)

	account := receiver.Accounts[0].Account.(*protocol.VMessAccount)
	assert.String(account.ID.String()).Equals("e641f5ad-9397-41e3-bf1a-e8740dfed019")
}
