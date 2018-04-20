package buf_test

import (
	"bytes"
	"io"
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/transport/pipe"
	. "v2ray.com/ext/assert"
)

func TestAdaptiveReader(t *testing.T) {
	assert := With(t)

	reader := NewReader(bytes.NewReader(make([]byte, 1024*1024)))
	b, err := reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b.Len(), Equals, int32(2*1024))

	b, err = reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b.Len(), Equals, int32(8*1024))

	b, err = reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b.Len(), Equals, int32(32*1024))

	b, err = reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b.Len(), Equals, int32(128*1024))

	b, err = reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(b.Len(), Equals, int32(512*1024))
}

func TestBytesReaderWriteTo(t *testing.T) {
	assert := With(t)

	pReader, pWriter := pipe.New()
	reader := &BufferedReader{Reader: pReader}
	b1 := New()
	b1.AppendBytes('a', 'b', 'c')
	b2 := New()
	b2.AppendBytes('e', 'f', 'g')
	assert(pWriter.WriteMultiBuffer(NewMultiBufferValue(b1, b2)), IsNil)
	pWriter.Close()

	pReader2, pWriter2 := pipe.New()
	writer := NewBufferedWriter(pWriter2)
	writer.SetBuffered(false)

	nBytes, err := io.Copy(writer, reader)
	assert(err, IsNil)
	assert(nBytes, Equals, int64(6))

	mb, err := pReader2.ReadMultiBuffer()
	assert(err, IsNil)
	assert(len(mb), Equals, 2)
	assert(mb[0].String(), Equals, "abc")
	assert(mb[1].String(), Equals, "efg")
}

func TestBytesReaderMultiBuffer(t *testing.T) {
	assert := With(t)

	pReader, pWriter := pipe.New()
	reader := &BufferedReader{Reader: pReader}
	b1 := New()
	b1.AppendBytes('a', 'b', 'c')
	b2 := New()
	b2.AppendBytes('e', 'f', 'g')
	assert(pWriter.WriteMultiBuffer(NewMultiBufferValue(b1, b2)), IsNil)
	pWriter.Close()

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
