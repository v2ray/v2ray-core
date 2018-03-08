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

	writer := NewBufferedWriter(NewWriter(writeBuffer))
	writer.SetBuffered(false)
	err := writer.WriteMultiBuffer(NewMultiBufferValue(lb))
	assert(err, IsNil)
	assert(writer.Flush(), IsNil)
	assert(expectedBytes, Equals, writeBuffer.Bytes())
}

func TestBytesWriterReadFrom(t *testing.T) {
	assert := With(t)

	cache := ray.NewStream(context.Background())
	const size = 50000
	reader := bufio.NewReader(io.LimitReader(rand.Reader, size))
	writer := NewBufferedWriter(cache)
	writer.SetBuffered(false)
	nBytes, err := reader.WriteTo(writer)
	assert(nBytes, Equals, int64(size))
	assert(err, IsNil)

	mb, err := cache.ReadMultiBuffer()
	assert(err, IsNil)
	assert(mb.Len(), Equals, size)
}

func TestDiscardBytes(t *testing.T) {
	assert := With(t)

	b := New()
	common.Must(b.Reset(ReadFullFrom(rand.Reader, Size)))

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
	nBytes, err := io.Copy(DiscardBytes, NewBufferedReader(r))
	assert(nBytes, Equals, int64(size))
	assert(err, IsNil)
}

func TestWriterInterface(t *testing.T) {
	assert := With(t)

	assert((*BufferToBytesWriter)(nil), Implements, (*Writer)(nil))
	assert((*BufferToBytesWriter)(nil), Implements, (*io.Writer)(nil))
	assert((*BufferToBytesWriter)(nil), Implements, (*io.ReaderFrom)(nil))

	assert((*BufferedWriter)(nil), Implements, (*Writer)(nil))
	assert((*BufferedWriter)(nil), Implements, (*io.Writer)(nil))
	assert((*BufferedWriter)(nil), Implements, (*io.ReaderFrom)(nil))
	assert((*BufferedWriter)(nil), Implements, (*io.ByteWriter)(nil))
}
