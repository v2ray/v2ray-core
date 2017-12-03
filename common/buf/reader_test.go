package buf_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/transport/ray"
	. "v2ray.com/ext/assert"
)

func TestAdaptiveReader(t *testing.T) {
	assert := With(t)

	reader := NewReader(bytes.NewReader(make([]byte, 1024*1024)))
	b, err := reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b.Len(), Equals, 2*1024)

	b, err = reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b.Len(), Equals, 32*1024)
}

func TestBytesReaderWriteTo(t *testing.T) {
	assert := With(t)

	stream := ray.NewStream(context.Background())
	reader := NewBufferedReader(stream)
	b1 := New()
	b1.AppendBytes('a', 'b', 'c')
	b2 := New()
	b2.AppendBytes('e', 'f', 'g')
	assert(stream.WriteMultiBuffer(NewMultiBufferValue(b1, b2)), IsNil)
	stream.Close()

	stream2 := ray.NewStream(context.Background())
	writer := NewBufferedWriter(stream2)
	writer.SetBuffered(false)

	nBytes, err := io.Copy(writer, reader)
	assert(err, IsNil)
	assert(nBytes, Equals, int64(6))

	mb, err := stream2.ReadMultiBuffer()
	assert(err, IsNil)
	assert(len(mb), Equals, 2)
	assert(mb[0].String(), Equals, "abc")
	assert(mb[1].String(), Equals, "efg")
}

func TestBytesReaderMultiBuffer(t *testing.T) {
	assert := With(t)

	stream := ray.NewStream(context.Background())
	reader := NewBufferedReader(stream)
	b1 := New()
	b1.AppendBytes('a', 'b', 'c')
	b2 := New()
	b2.AppendBytes('e', 'f', 'g')
	assert(stream.WriteMultiBuffer(NewMultiBufferValue(b1, b2)), IsNil)
	stream.Close()

	mbReader := NewReader(reader)
	mb, err := mbReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(len(mb), Equals, 2)
	assert(mb[0].String(), Equals, "abc")
	assert(mb[1].String(), Equals, "efg")
}

func TestReaderInterface(t *testing.T) {
	assert := With(t)

	assert((*BytesToBufferReader)(nil), Implements, (*io.Reader)(nil))
	assert((*BytesToBufferReader)(nil), Implements, (*Reader)(nil))

	assert((*BufferedReader)(nil), Implements, (*Reader)(nil))
	assert((*BufferedReader)(nil), Implements, (*io.Reader)(nil))
	assert((*BufferedReader)(nil), Implements, (*io.ByteReader)(nil))
	assert((*BufferedReader)(nil), Implements, (*io.WriterTo)(nil))
}
