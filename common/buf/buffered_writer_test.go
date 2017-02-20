package buf_test

import (
	"crypto/rand"
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
)

func TestBufferedWriter(t *testing.T) {
	assert := assert.On(t)

	content := New()

	writer := NewBufferedWriter(content)
	assert.Bool(writer.IsBuffered()).IsTrue()

	payload := make([]byte, 16)

	nBytes, err := writer.Write(payload)
	assert.Int(nBytes).Equals(16)
	assert.Error(err).IsNil()

	assert.Bool(content.IsEmpty()).IsTrue()

	assert.Error(writer.SetBuffered(false)).IsNil()
	assert.Int(content.Len()).Equals(16)
}

func TestBufferedWriterLargePayload(t *testing.T) {
	assert := assert.On(t)

	content := NewLocal(128 * 1024)

	writer := NewBufferedWriter(content)
	assert.Bool(writer.IsBuffered()).IsTrue()

	payload := make([]byte, 64*1024)
	rand.Read(payload)

	nBytes, err := writer.Write(payload[:512])
	assert.Int(nBytes).Equals(512)
	assert.Error(err).IsNil()

	assert.Bool(content.IsEmpty()).IsTrue()

	nBytes, err = writer.Write(payload[512:])
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(64*1024 - 512)
	assert.Bytes(content.Bytes()).Equals(payload)
}
