package serial_test

import (
	"testing"

	. "v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
)

func TestGetInstance(t *testing.T) {
	assert := assert.On(t)

	p, err := GetInstance("")
	assert.Pointer(p).IsNil()
	assert.Error(err).IsNotNil()
}
