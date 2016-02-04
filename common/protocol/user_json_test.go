// +build json

package protocol_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/protocol"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestUserParsing(t *testing.T) {
	v2testing.Current(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{
    "id": "96edb838-6d68-42ef-a933-25f7ac3a9d09",
    "email": "love@v2ray.com",
    "level": 1,
    "alterId": 100
  }`), user)
	assert.Error(err).IsNil()
	assert.String(user.ID).Equals("96edb838-6d68-42ef-a933-25f7ac3a9d09")
	assert.Byte(byte(user.Level)).Equals(1)
}

func TestInvalidUserJson(t *testing.T) {
	v2testing.Current(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{"id": 1234}`), user)
	assert.Error(err).IsNotNil()
}

func TestInvalidIdJson(t *testing.T) {
	v2testing.Current(t)

	user := new(User)
	err := json.Unmarshal([]byte(`{"id": "1234"}`), user)
	assert.Error(err).IsNotNil()
}
