package socks_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	_ "v2ray.com/core/common/net/testing"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/core/proxy/socks"
	. "v2ray.com/ext/assert"
)

func TestUDPEncoding(t *testing.T) {
	assert := With(t)

	b := buf.New()

	request := &protocol.RequestHeader{
		Address: net.IPAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}),
		Port:    1024,
	}
	writer := buf.NewSequentialWriter(NewUDPWriter(request, b))

	content := []byte{'a'}
	payload := buf.New()
	payload.Write(content)
	assert(writer.WriteMultiBuffer(buf.NewMultiBufferValue(payload)), IsNil)

	reader := NewUDPReader(b)

	decodedPayload, err := reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(decodedPayload[0].Bytes(), Equals, content)
}
