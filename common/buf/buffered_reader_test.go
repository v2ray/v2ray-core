package buf_test

import (
	"crypto/rand"
	"testing"

	. "v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
)

func TestBufferedReader(t *testing.T) {
	assert := With(t)

	content := New()
	assert(content.AppendSupplier(ReadFrom(rand.Reader)), IsNil)

	len := content.Len()

	reader := NewBufferedReader(content)
	assert(reader.IsBuffered(), IsTrue)

	payload := make([]byte, 16)

	nBytes, err := reader.Read(payload)
	assert(nBytes, Equals, 16)
	assert(err, IsNil)

	len2 := content.Len()
	assert(len - len2, GreaterThan, 16)

	nBytes, err = reader.Read(payload)
	assert(nBytes, Equals, 16)
	assert(err, IsNil)

	assert(content.Len(), Equals, len2)
}
