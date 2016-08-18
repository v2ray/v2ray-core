// +build json

package collect_test

import (
	"encoding/json"
	"testing"

	. "github.com/v2ray/v2ray-core/common/collect"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestStringListUnmarshalError(t *testing.T) {
	assert := assert.On(t)

	rawJson := `1234`
	list := new(StringList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert.Error(err).IsNotNil()
}

func TestStringListLen(t *testing.T) {
	assert := assert.On(t)

	rawJson := `"a, b, c, d"`
	list := new(StringList)
	err := json.Unmarshal([]byte(rawJson), list)
	assert.Error(err).IsNil()
	assert.Int(list.Len()).Equals(4)
}
