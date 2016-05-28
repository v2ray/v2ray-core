// +build json

package protocol_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/testing/assert"
)

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
	assert.Byte(byte(user.Level)).Equals(1)

	account, ok := user.Account.(*VMessAccount)
	assert.Bool(ok).IsTrue()
	assert.String(account.ID.String()).Equals("96edb838-6d68-42ef-a933-25f7ac3a9d09")
}

func TestInvalidUserJson(t *testing.T) {
	assert := assert.On(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{"id": 1234}`), user)
	assert.Error(err).IsNotNil()
}

func TestInvalidIdJson(t *testing.T) {
	assert := assert.On(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{"id": "1234"}`), user)
	assert.Error(err).IsNotNil()
}
