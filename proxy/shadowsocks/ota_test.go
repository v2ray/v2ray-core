package shadowsocks_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/proxy/shadowsocks"
	"v2ray.com/core/testing/assert"
)

func TestNormalChunkReading(t *testing.T) {
	assert := assert.On(t)

	buffer := buf.New()
	buffer.AppendBytes(
		0, 8, 39, 228, 69, 96, 133, 39, 254, 26, 201, 70, 11, 12, 13, 14, 15, 16, 17, 18)
	reader := NewChunkReader(buffer, NewAuthenticator(ChunkKeyGenerator(
		[]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36})))
	payload, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bytes(payload[0].Bytes()).Equals([]byte{11, 12, 13, 14, 15, 16, 17, 18})
}

func TestNormalChunkWriting(t *testing.T) {
	assert := assert.On(t)

	buffer := buf.NewLocal(512)
	writer := NewChunkWriter(buffer, NewAuthenticator(ChunkKeyGenerator(
		[]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36})))

	b := buf.NewLocal(256)
	b.Append([]byte{11, 12, 13, 14, 15, 16, 17, 18})
	err := writer.Write(buf.NewMultiBufferValue(b))
	assert.Error(err).IsNil()
	assert.Bytes(buffer.Bytes()).Equals([]byte{0, 8, 39, 228, 69, 96, 133, 39, 254, 26, 201, 70, 11, 12, 13, 14, 15, 16, 17, 18})
}
