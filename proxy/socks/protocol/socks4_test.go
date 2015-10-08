package protocol

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestSocks4AuthenticationRequestRead(t *testing.T) {
	assert := unit.Assert(t)

	rawRequest := []byte{
		0x04, // version
		0x01, // command
		0x00, 0x35,
		0x72, 0x72, 0x72, 0x72,
	}
	_, request4, err := ReadAuthentication(bytes.NewReader(rawRequest))
	assert.Error(err).HasCode(1000)
	assert.Byte(request4.Version).Named("Version").Equals(0x04)
	assert.Byte(request4.Command).Named("Command").Equals(0x01)
	assert.Uint16(request4.Port).Named("Port").Equals(53)
	assert.Bytes(request4.IP[:]).Named("IP").Equals([]byte{0x72, 0x72, 0x72, 0x72})
}

func TestSocks4AuthenticationResponseToBytes(t *testing.T) {
	assert := unit.Assert(t)

	response := &Socks4AuthenticationResponse{
		result: byte(0x10),
		port:   443,
		ip:     []byte{1, 2, 3, 4},
	}

	buffer := alloc.NewSmallBuffer().Clear()
	defer buffer.Release()

	response.Write(buffer)
	assert.Bytes(buffer.Value).Equals([]byte{0x00, 0x10, 0x01, 0xBB, 0x01, 0x02, 0x03, 0x04})
}
