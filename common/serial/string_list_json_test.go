// +build json

package serial_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestStringListUnmarshalError(t *testing.T) {
	assert := assert.On(t)

	rawJson := `1234`
	list := new(StringTList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert.Error(err).IsNotNil()
}
