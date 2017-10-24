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

	rawContent := make([]byte, 1024*1024)
	buffer := bytes.NewBuffer(rawContent)

	reader := NewReader(buffer)
	b, err := reader.Read()
	assert(err, IsNil)
	assert(b.Len(), Equals, 32*1024)
}

func TestBytesReaderWriteTo(t *testing.T) {
	assert := With(t)

	stream := ray.NewStream(context.Background())
	reader := ToBytesReader(stream)
	b1 := New()
	b1.AppendBytes('a', 'b', 'c')
	b2 := New()
	b2.AppendBytes('e', 'f', 'g')
	assert(stream.Write(NewMultiBufferValue(b1, b2)), IsNil)
	stream.Close()

	stream2 := ray.NewStream(context.Background())
	writer := ToBytesWriter(stream2)

	nBytes, err := io.Copy(writer, reader)
	assert(err, IsNil)
	assert(nBytes, Equals, int64(6))

	mb, err := stream2.Read()
	assert(err, IsNil)
	assert(len(mb), Equals, 2)
	assert(mb[0].String(), Equals, "abc")
	assert(mb[1].String(), Equals, "efg")
}

func TestBytesReaderMultiBuffer(t *testing.T) {
	assert := With(t)

	stream := ray.NewStream(context.Background())
	reader := ToBytesReader(stream)
	b1 := New()
	b1.AppendBytes('a', 'b', 'c')
	b2 := New()
	b2.AppendBytes('e', 'f', 'g')
	assert(stream.Write(NewMultiBufferValue(b1, b2)), IsNil)
	stream.Close()

	mbReader := NewReader(reader)
	mb, err := mbReader.Read()
	assert(err, IsNil)
	assert(len(mb), Equals, 2)
	assert(mb[0].String(), Equals, "abc")
	assert(mb[1].String(), Equals, "efg")
}
