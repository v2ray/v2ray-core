package ray_test

import (
	"context"
	"io"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/ray"
	. "v2ray.com/ext/assert"
)

func TestStreamIO(t *testing.T) {
	assert := With(t)

	stream := NewStream(context.Background())
	b1 := buf.New()
	b1.AppendBytes('a')
	assert(stream.WriteMultiBuffer(buf.NewMultiBufferValue(b1)), IsNil)

	_, err := stream.ReadMultiBuffer()
	assert(err, IsNil)

	stream.Close()
	_, err = stream.ReadMultiBuffer()
	assert(err, Equals, io.EOF)

	b2 := buf.New()
	b2.AppendBytes('b')
	err = stream.WriteMultiBuffer(buf.NewMultiBufferValue(b2))
	assert(err, Equals, io.ErrClosedPipe)
}

func TestStreamClose(t *testing.T) {
	assert := With(t)

	stream := NewStream(context.Background())
	b1 := buf.New()
	b1.AppendBytes('a')
	assert(stream.WriteMultiBuffer(buf.NewMultiBufferValue(b1)), IsNil)

	stream.Close()

	_, err := stream.ReadMultiBuffer()
	assert(err, IsNil)

	_, err = stream.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
}
