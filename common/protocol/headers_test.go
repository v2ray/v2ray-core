package protocol_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestRequestOptionSet(t *testing.T) {
	assert := assert.On(t)

	option := new(RequestOption)
	assert.Bool(option.Has(RequestOptionChunkStream)).IsFalse()

	option.Set(RequestOptionChunkStream)
	assert.Bool(option.Has(RequestOptionChunkStream)).IsTrue()

	option.Set(RequestOptionConnectionReuse)
	assert.Bool(option.Has(RequestOptionConnectionReuse)).IsTrue()
	assert.Bool(option.Has(RequestOptionChunkStream)).IsTrue()
}

func TestRequestOptionClear(t *testing.T) {
	assert := assert.On(t)

	option := new(RequestOption)
	option.Set(RequestOptionChunkStream)
	option.Set(RequestOptionConnectionReuse)

	option.Clear(RequestOptionChunkStream)
	assert.Bool(option.Has(RequestOptionChunkStream)).IsFalse()
	assert.Bool(option.Has(RequestOptionConnectionReuse)).IsTrue()
}
