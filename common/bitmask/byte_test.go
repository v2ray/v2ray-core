package bitmask_test

import (
	"testing"

	. "v2ray.com/core/common/bitmask"
	"v2ray.com/core/testing/assert"
)

func TestBitmaskByte(t *testing.T) {
	assert := assert.On(t)

	b := Byte(0)
	b.Set(Byte(1))
	assert.Bool(b.Has(1)).IsTrue()

	b.Set(Byte(2))
	assert.Bool(b.Has(2)).IsTrue()
	assert.Bool(b.Has(1)).IsTrue()

	b.Clear(Byte(1))
	assert.Bool(b.Has(2)).IsTrue()
	assert.Bool(b.Has(1)).IsFalse()

	b.Toggle(Byte(2))
	assert.Bool(b.Has(2)).IsFalse()
}
