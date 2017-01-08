package socks_test

import (
	"testing"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/core/proxy/socks"
	"v2ray.com/core/testing/assert"
)

func TestUDPEncoding(t *testing.T) {
	assert := assert.On(t)

	b := buf.New()

	request := &protocol.RequestHeader{
		Address: net.IPAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6}),
		Port:    1024,
	}
	writer := NewUDPWriter(request, b)

	content := []byte{'a'}
	payload := buf.New()
	payload.Append(content)
	assert.Error(writer.Write(payload)).IsNil()

	reader := NewUDPReader(b)

	decodedPayload, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bytes(decodedPayload.Bytes()).Equals(content)
}
