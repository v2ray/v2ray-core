package pipe_test

import (
	"io"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/pipe"
	. "v2ray.com/ext/assert"
)

func TestPipeReadWrite(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New()
	payload := []byte{'a', 'b', 'c', 'd'}
	b := buf.New()
	b.Write(payload)
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)

	rb, err := pReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(rb.String(), Equals, b.String())
}

func TestPipeCloseError(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New()
	payload := []byte{'a', 'b', 'c', 'd'}
	b := buf.New()
	b.Write(payload)
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)
	pWriter.CloseError()

	rb, err := pReader.ReadMultiBuffer()
	assert(err, Equals, io.ErrClosedPipe)
	assert(rb.IsEmpty(), IsTrue)
}

func TestPipeClose(t *testing.T) {
	assert := With(t)

	pReader, pWriter := New()
	payload := []byte{'a', 'b', 'c', 'd'}
	b := buf.New()
	b.Write(payload)
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)
	assert(pWriter.Close(), IsNil)

	rb, err := pReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(rb.String(), Equals, b.String())

	rb, err = pReader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
	assert(rb.IsEmpty(), IsTrue)
}
