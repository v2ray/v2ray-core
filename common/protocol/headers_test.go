package protocol_test

import (
	"testing"

	. "v2ray.com/core/common/protocol"
	"v2ray.com/core/testing/assert"
)

func TestRequestOptionSet(t *testing.T) {
	assert := assert.On(t)

	var option RequestOption
	assert.Bool(option.Has(RequestOptionChunkStream)).IsFalse()

	option.Set(RequestOptionChunkStream)
	assert.Bool(option.Has(RequestOptionChunkStream)).IsTrue()

	option.Set(RequestOptionChunkMasking)
	assert.Bool(option.Has(RequestOptionChunkMasking)).IsTrue()
	assert.Bool(option.Has(RequestOptionChunkStream)).IsTrue()
}

func TestRequestOptionClear(t *testing.T) {
	assert := assert.On(t)

	var option RequestOption
	option.Set(RequestOptionChunkStream)
	option.Set(RequestOptionChunkMasking)

	option.Clear(RequestOptionChunkStream)
	assert.Bool(option.Has(RequestOptionChunkStream)).IsFalse()
	assert.Bool(option.Has(RequestOptionChunkMasking)).IsTrue()
}
