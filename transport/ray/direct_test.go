package ray_test

import (
	"io"
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/ray"
)

func TestStreamIO(t *testing.T) {
	assert := assert.On(t)

	stream := NewStream()
	assert.Error(stream.Write(buf.New())).IsNil()

	_, err := stream.Read()
	assert.Error(err).IsNil()

	stream.Close()
	_, err = stream.Read()
	assert.Error(err).Equals(io.EOF)

	err = stream.Write(buf.New())
	assert.Error(err).Equals(io.ErrClosedPipe)
}

func TestStreamClose(t *testing.T) {
	assert := assert.On(t)

	stream := NewStream()
	assert.Error(stream.Write(buf.New())).IsNil()

	stream.Close()

	_, err := stream.Read()
	assert.Error(err).IsNil()

	_, err = stream.Read()
	assert.Error(err).Equals(io.EOF)
}
