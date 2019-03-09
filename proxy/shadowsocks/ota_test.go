package shadowsocks_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	. "v2ray.com/core/proxy/shadowsocks"
)

func TestNormalChunkReading(t *testing.T) {
	buffer := buf.New()
	buffer.Write([]byte{0, 8, 39, 228, 69, 96, 133, 39, 254, 26, 201, 70, 11, 12, 13, 14, 15, 16, 17, 18})
	reader := NewChunkReader(buffer, NewAuthenticator(ChunkKeyGenerator(
		[]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36})))
	payload, err := reader.ReadMultiBuffer()
	common.Must(err)

	if diff := cmp.Diff(payload[0].Bytes(), []byte{11, 12, 13, 14, 15, 16, 17, 18}); diff != "" {
		t.Error(diff)
	}
}

func TestNormalChunkWriting(t *testing.T) {
	buffer := buf.New()
	writer := NewChunkWriter(buffer, NewAuthenticator(ChunkKeyGenerator(
		[]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36})))

	b := buf.New()
	b.Write([]byte{11, 12, 13, 14, 15, 16, 17, 18})
	common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{b}))
	if diff := cmp.Diff(buffer.Bytes(), []byte{0, 8, 39, 228, 69, 96, 133, 39, 254, 26, 201, 70, 11, 12, 13, 14, 15, 16, 17, 18}); diff != "" {
		t.Error(diff)
	}
}
