package buf_test

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"testing"

	"context"
	"io"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/ray"
)

func TestWriter(t *testing.T) {
	assert := assert.On(t)

	lb := New()
	assert.Error(lb.AppendSupplier(ReadFrom(rand.Reader))).IsNil()

	expectedBytes := append([]byte(nil), lb.Bytes()...)

	writeBuffer := bytes.NewBuffer(make([]byte, 0, 1024*1024))

	writer := NewWriter(NewBufferedWriter(writeBuffer))
	err := writer.Write(NewMultiBufferValue(lb))
	assert.Error(err).IsNil()
	assert.Bytes(expectedBytes).Equals(writeBuffer.Bytes())
}

func TestBytesWriterReadFrom(t *testing.T) {
	assert := assert.On(t)

	cache := ray.NewStream(context.Background())
	reader := bufio.NewReader(io.LimitReader(rand.Reader, 8192))
	_, err := reader.WriteTo(ToBytesWriter(cache))
	assert.Error(err).IsNil()

	mb, err := cache.Read()
	assert.Error(err).IsNil()
	assert.Int(mb.Len()).Equals(8192)
	assert.Int(len(mb)).Equals(4)
}
