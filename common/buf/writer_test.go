package buf_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
)

func TestWriter(t *testing.T) {
	assert := assert.On(t)

	lb := New()
	assert.Error(lb.AppendSupplier(ReadFrom(rand.Reader))).IsNil()

	expectedBytes := append([]byte(nil), lb.Bytes()...)

	writeBuffer := bytes.NewBuffer(make([]byte, 0, 1024*1024))

	writer := NewWriter(NewBufferedWriter(writeBuffer))
	err := writer.Write(lb)
	assert.Error(err).IsNil()
	assert.Bytes(expectedBytes).Equals(writeBuffer.Bytes())
}
