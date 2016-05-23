// +build json

package serial_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/serial"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestStringListUnmarshalError(t *testing.T) {
	v2testing.Current(t)

	rawJson := `1234`
	list := new(StringTList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert.Error(err).IsNotNil()
}
