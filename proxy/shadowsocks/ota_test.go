package shadowsocks_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/proxy/shadowsocks"
	. "v2ray.com/ext/assert"
)

func TestNormalChunkReading(t *testing.T) {
	assert := With(t)

	buffer := buf.New()
	buffer.WriteBytes(
		0, 8, 39, 228, 69, 96, 133, 39, 254, 26, 201, 70, 11, 12, 13, 14, 15, 16, 17, 18)
	reader := NewChunkReader(buffer, NewAuthenticator(ChunkKeyGenerator(
		[]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36})))
	payload, err := reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(payload[0].Bytes(), Equals, []byte{11, 12, 13, 14, 15, 16, 17, 18})
}

func TestNormalChunkWriting(t *testing.T) {
	assert := With(t)

	buffer := buf.NewSize(512)
	writer := NewChunkWriter(buffer, NewAuthenticator(ChunkKeyGenerator(
		[]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36})))

	b := buf.NewSize(256)
	b.Write([]byte{11, 12, 13, 14, 15, 16, 17, 18})
	err := writer.WriteMultiBuffer(buf.NewMultiBufferValue(b))
	assert(err, IsNil)
	assert(buffer.Bytes(), Equals, []byte{0, 8, 39, 228, 69, 96, 133, 39, 254, 26, 201, 70, 11, 12, 13, 14, 15, 16, 17, 18})
}
