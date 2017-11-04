package buf_test

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"testing"

	"context"
	"io"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/buf"
	"v2ray.com/core/transport/ray"
	. "v2ray.com/ext/assert"
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

func TestDiscardBytes(t *testing.T) {
	assert := With(t)

	b := New()
	common.Must(b.Reset(ReadFrom(rand.Reader)))

	nBytes, err := io.Copy(DiscardBytes, b)
	assert(nBytes, Equals, int64(Size))
	assert(err, IsNil)
}

func TestDiscardBytesMultiBuffer(t *testing.T) {
	assert := With(t)

	const size = 10240*1024 + 1
	buffer := bytes.NewBuffer(make([]byte, 0, size))
	common.Must2(buffer.ReadFrom(io.LimitReader(rand.Reader, size)))

	r := NewReader(buffer)
	nBytes, err := io.Copy(DiscardBytes, ToBytesReader(r))
	assert(nBytes, Equals, int64(size))
	assert(err, IsNil)
}
