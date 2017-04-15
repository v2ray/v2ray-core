package buf_test

import (
	"testing"

	"context"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/ray"
)

func TestMergingReader(t *testing.T) {
	assert := assert.On(t)

	stream := ray.NewStream(context.Background())
	b1 := New()
	b1.AppendBytes('a', 'b', 'c')
	stream.Write(NewMultiBufferValue(b1))

	b2 := New()
	b2.AppendBytes('e', 'f', 'g')
	stream.Write(NewMultiBufferValue(b2))

	b3 := New()
	b3.AppendBytes('h', 'i', 'j')
	stream.Write(NewMultiBufferValue(b3))

	reader := NewMergingReader(stream)
	b, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Int(b.Len()).Equals(9)
}
