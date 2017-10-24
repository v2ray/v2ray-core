package buf_test

import (
	"crypto/rand"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
)

func TestBufferedWriter(t *testing.T) {
	assert := With(t)

	content := New()

	writer := NewBufferedWriter(content)
	assert(writer.IsBuffered(), IsTrue)

	payload := make([]byte, 16)

	nBytes, err := writer.Write(payload)
	assert(nBytes, Equals, 16)
	assert(err, IsNil)

	assert(content.IsEmpty(), IsTrue)

	assert(writer.SetBuffered(false), IsNil)
	assert(content.Len(), Equals, 16)
}

func TestBufferedWriterLargePayload(t *testing.T) {
	assert := With(t)

	content := NewLocal(128 * 1024)

	writer := NewBufferedWriter(content)
	assert(writer.IsBuffered(), IsTrue)

	payload := make([]byte, 64*1024)
	common.Must2(rand.Read(payload))

	nBytes, err := writer.Write(payload[:512])
	assert(nBytes, Equals, 512)
	assert(err, IsNil)

	assert(content.IsEmpty(), IsTrue)

	nBytes, err = writer.Write(payload[512:])
	assert(err, IsNil)
	assert(writer.Flush(), IsNil)
	assert(nBytes, Equals, 64*1024 - 512)
	assert(content.Bytes(), Equals, payload)
}
