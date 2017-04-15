package buf_test

import (
	"bytes"
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
)

func TestAdaptiveReader(t *testing.T) {
	assert := assert.On(t)

	rawContent := make([]byte, 1024*1024)
	buffer := bytes.NewBuffer(rawContent)

	reader := NewReader(buffer)
	b, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Int(b.Len()).Equals(32 * 1024)
}
