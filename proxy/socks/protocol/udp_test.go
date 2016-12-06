package protocol

import (
	"testing"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func TestSingleByteUDPRequest(t *testing.T) {
	assert := assert.On(t)

	request, err := ReadUDPRequest(make([]byte, 1))
	if request != nil {
		t.Fail()
	}
	assert.Error(err).IsNotNil()
}

func TestDomainAddressRequest(t *testing.T) {
	assert := assert.On(t)

	payload := make([]byte, 0, 1024)
	payload = append(payload, 0, 0, 1, AddrTypeDomain, byte(len("v2ray.com")))
	payload = append(payload, []byte("v2ray.com")...)
	payload = append(payload, 0, 80)
	payload = append(payload, []byte("Actual payload")...)

	request, err := ReadUDPRequest(payload)
	assert.Error(err).IsNil()

	assert.Byte(request.Fragment).Equals(1)
	assert.Address(request.Address).EqualsString("v2ray.com")
	assert.Port(request.Port).Equals(v2net.Port(80))
	assert.String(request.Data.String()).Equals("Actual payload")
}
