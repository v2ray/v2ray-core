// +build json

package protocol_test

import (
	"encoding/json"
	"testing"

	. "v2ray.com/core/common/protocol"
	"v2ray.com/core/testing/assert"
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
	assert.String(user.Email).Equals("love@v2ray.com")
}

func TestInvalidUserJson(t *testing.T) {
	assert := assert.On(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{"email": 1234}`), user)
	assert.Error(err).IsNotNil()
}
