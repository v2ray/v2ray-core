package shadowsocks_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	. "v2ray.com/core/proxy/shadowsocks"
	"v2ray.com/core/testing/assert"
)

func TestNormalChunkReading(t *testing.T) {
	assert := assert.On(t)

	buffer := alloc.NewBuffer().Clear().AppendBytes(
		0, 8, 39, 228, 69, 96, 133, 39, 254, 26, 201, 70, 11, 12, 13, 14, 15, 16, 17, 18)
	reader := NewChunkReader(buffer, NewAuthenticator(ChunkKeyGenerator(
		[]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36})))
	payload, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bytes(payload.Value).Equals([]byte{11, 12, 13, 14, 15, 16, 17, 18})

	payload.PrependBytes(3, 4)
	assert.Bytes(payload.Value).Equals([]byte{3, 4, 11, 12, 13, 14, 15, 16, 17, 18})
}
