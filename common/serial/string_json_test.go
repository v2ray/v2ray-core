// +build json

package serial_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestInvalidStringTJson(t *testing.T) {
	assert := assert.On(t)

	var s StringT
	err := json.Unmarshal([]byte("1"), &s)
	assert.Error(err).IsNotNil()
}

func TestStringTParsing(t *testing.T) {
	assert := assert.On(t)

	var s StringT
	err := json.Unmarshal([]byte("\"1\""), &s)
	assert.Error(err).IsNil()
	assert.String(s.String()).Equals("1")
}
