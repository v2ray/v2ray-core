// +build json

package serial_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/serial"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestInvalidStringLiteralJson(t *testing.T) {
	v2testing.Current(t)

	var s StringLiteral
	err := json.Unmarshal([]byte("1"), &s)
	assert.Error(err).IsNotNil()
}

func TestStringLiteralParsing(t *testing.T) {
	v2testing.Current(t)

	var s StringLiteral
	err := json.Unmarshal([]byte("\"1\""), &s)
	assert.Error(err).IsNil()
	assert.String(s).Equals("1")
}
