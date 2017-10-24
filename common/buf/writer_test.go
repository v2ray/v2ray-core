package buf_test

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"testing"

	"context"
	"io"

	. "v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
	"v2ray.com/core/transport/ray"
)

func TestWriter(t *testing.T) {
	assert := With(t)

	lb := New()
	assert(lb.AppendSupplier(ReadFrom(rand.Reader)), IsNil)

	expectedBytes := append([]byte(nil), lb.Bytes()...)

	writeBuffer := bytes.NewBuffer(make([]byte, 0, 1024*1024))

	writer := NewWriter(NewBufferedWriter(writeBuffer))
	err := writer.Write(NewMultiBufferValue(lb))
	assert(err, IsNil)
	assert(expectedBytes, Equals, writeBuffer.Bytes())
}

func TestBytesWriterReadFrom(t *testing.T) {
	assert := With(t)

	cache := ray.NewStream(context.Background())
	reader := bufio.NewReader(io.LimitReader(rand.Reader, 8192))
	_, err := reader.WriteTo(ToBytesWriter(cache))
	assert(err, IsNil)

	mb, err := cache.Read()
	assert(err, IsNil)
	assert(mb.Len(), Equals, 8192)
	assert(len(mb), Equals, 4)
}
