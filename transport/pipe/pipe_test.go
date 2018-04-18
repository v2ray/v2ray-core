package pipe_test

import (
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
	b.Append(payload)
	assert(pWriter.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)

	rb, err := pReader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(rb.String(), Equals, b.String())
}
