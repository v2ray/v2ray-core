package protocol

import (
	"bytes"
	"testing"

	"v2ray.com/core/common/alloc"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestSocks4AuthenticationRequestRead(t *testing.T) {
	assert := assert.On(t)

	rawRequest := []byte{
		0x04, // version
		0x01, // command
		0x00, 0x35,
		0x72, 0x72, 0x72, 0x72,
	}
	_, request4, err := ReadAuthentication(bytes.NewReader(rawRequest))
	assert.Error(err).Equals(Socks4Downgrade)
	assert.Byte(request4.Version).Equals(0x04)
	assert.Byte(request4.Command).Equals(0x01)
	assert.Port(request4.Port).Equals(v2net.Port(53))
	assert.Bytes(request4.IP[:]).Equals([]byte{0x72, 0x72, 0x72, 0x72})
}

func TestSocks4AuthenticationResponseToBytes(t *testing.T) {
	assert := assert.On(t)

	response := NewSocks4AuthenticationResponse(byte(0x10), 443, []byte{1, 2, 3, 4})

	buffer := alloc.NewLocalBuffer(2048)
	defer buffer.Release()

	response.Write(buffer)
	assert.Bytes(buffer.Bytes()).Equals([]byte{0x00, 0x10, 0x01, 0xBB, 0x01, 0x02, 0x03, 0x04})
}
