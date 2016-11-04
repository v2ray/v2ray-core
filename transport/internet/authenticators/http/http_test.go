package http_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/authenticators/http"
)

func TestReaderWriter(t *testing.T) {
	assert := assert.On(t)

	cache := alloc.NewBuffer()
	writer := NewHeaderWriter(alloc.NewLocalBuffer(256).Clear().AppendString("abcd" + ENDING))
	writer.Write(cache)
	cache.Write([]byte{'e', 'f', 'g'})

	reader := &HeaderReader{}
	buffer, err := reader.Read(cache)
	assert.Error(err).IsNil()
	assert.Bytes(buffer.Value).Equals([]byte{'e', 'f', 'g'})
}
