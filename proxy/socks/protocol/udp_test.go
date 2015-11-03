package protocol

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
	"github.com/v2ray/v2ray-core/transport"
)

func TestSingleByteRequest(t *testing.T) {
	assert := unit.Assert(t)

	request, err := ReadUDPRequest(make([]byte, 1))
	if request != nil {
		t.Fail()
	}
	assert.Error(err).Equals(transport.CorruptedPacket)
}

func TestDomainAddressRequest(t *testing.T) {
	assert := unit.Assert(t)

	payload := make([]byte, 0, 1024)
	payload = append(payload, 0, 0, 1, AddrTypeDomain, byte(len("v2ray.com")))
	payload = append(payload, []byte("v2ray.com")...)
	payload = append(payload, 0, 80)
	payload = append(payload, []byte("Actual payload")...)

	request, err := ReadUDPRequest(payload)
	assert.Error(err).IsNil()

	assert.Byte(request.Fragment).Equals(1)
	assert.String(request.Address.String()).Equals("v2ray.com:80")
	assert.Bytes(request.Data.Value).Equals([]byte("Actual payload"))
}
