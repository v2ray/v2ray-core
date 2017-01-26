package ray_test

import (
	"io"
	"testing"

	"context"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/ray"
)

func TestStreamIO(t *testing.T) {
	assert := assert.On(t)

	stream := NewStream(context.Background())
	b1 := buf.New()
	b1.AppendBytes('a')
	assert.Error(stream.Write(b1)).IsNil()

	_, err := stream.Read()
	assert.Error(err).IsNil()

	stream.Close()
	_, err = stream.Read()
	assert.Error(err).Equals(io.EOF)

	b2 := buf.New()
	b2.AppendBytes('b')
	err = stream.Write(b2)
	assert.Error(err).Equals(io.ErrClosedPipe)
}

func TestStreamClose(t *testing.T) {
	assert := assert.On(t)

	stream := NewStream(context.Background())
	b1 := buf.New()
	b1.AppendBytes('a')
	assert.Error(stream.Write(b1)).IsNil()

	stream.Close()

	_, err := stream.Read()
	assert.Error(err).IsNil()

	_, err = stream.Read()
	assert.Error(err).Equals(io.EOF)
}
