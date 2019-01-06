package serial_test

import (
	"testing"

	. "v2ray.com/core/common/serial"
	. "v2ray.com/ext/assert"
)

func TestGetInstance(t *testing.T) {
	assert := With(t)

	p, err := GetInstance("")
	assert(p, IsNil)
	assert(err, IsNotNil)
}

func TestConvertingNilMessage(t *testing.T) {
	x := ToTypedMessage(nil)
	if x != nil {
		t.Error("expect nil, but actually not")
	}
}
