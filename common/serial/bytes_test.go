package serial_test

import (
	"testing"

	. "v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
)

func TestBytesToHex(t *testing.T) {
	assert := assert.On(t)

	cases := []struct {
		input  []byte
		output string
	}{
		{input: []byte{}, output: "[]"},
		{input: []byte("a"), output: "[61]"},
		{input: []byte("abcd"), output: "[61,62,63,64]"},
		{input: []byte(";kdfpa;dfkaepr3ira;dlkvn;vopaehra;dkhf"), output: "[3b,6b,64,66,70,61,3b,64,66,6b,61,65,70,72,33,69,72,61,3b,64,6c,6b,76,6e,3b,76,6f,70,61,65,68,72,61,3b,64,6b,68,66]"},
	}

	for _, test := range cases {
		assert.String(test.output).Equals(BytesToHexString(test.input))
	}
}
