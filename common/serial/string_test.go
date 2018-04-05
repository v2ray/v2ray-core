package serial_test

import (
	"errors"
	"testing"

	. "v2ray.com/core/common/serial"
	. "v2ray.com/ext/assert"
)

func TestToString(t *testing.T) {
	assert := With(t)

	s := "a"
	data := []struct {
		Value  interface{}
		String string
	}{
		{Value: s, String: s},
		{Value: &s, String: s},
		{Value: errors.New("t"), String: "t"},
		{Value: []byte{'b', 'c'}, String: "[62,63]"},
	}

	for _, c := range data {
		assert(ToString(c.Value), Equals, c.String)
	}
}
