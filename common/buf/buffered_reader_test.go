package buf_test

import (
	"crypto/rand"
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
)

func TestBufferedReader(t *testing.T) {
	assert := assert.On(t)

	content := New()
	assert.Error(content.AppendSupplier(ReadFrom(rand.Reader))).IsNil()

	len := content.Len()

	reader := NewBufferedReader(content)
	assert.Bool(reader.IsBuffered()).IsTrue()

	payload := make([]byte, 16)

	nBytes, err := reader.Read(payload)
	assert.Int(nBytes).Equals(16)
	assert.Error(err).IsNil()

	len2 := content.Len()
	assert.Int(len - len2).GreaterThan(16)

	nBytes, err = reader.Read(payload)
	assert.Int(nBytes).Equals(16)
	assert.Error(err).IsNil()

	assert.Int(content.Len()).Equals(len2)
}
