package protocol

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/transport"
)

func TestSingleByteUDPRequest(t *testing.T) {
	v2testing.Current(t)

	request, err := ReadUDPRequest(make([]byte, 1))
	if request != nil {
		t.Fail()
	}
	assert.Error(err).Equals(transport.ErrorCorruptedPacket)
}

func TestDomainAddressRequest(t *testing.T) {
	v2testing.Current(t)

	payload := make([]byte, 0, 1024)
	payload = append(payload, 0, 0, 1, AddrTypeDomain, byte(len("v2ray.com")))
	payload = append(payload, []byte("v2ray.com")...)
	payload = append(payload, 0, 80)
	payload = append(payload, []byte("Actual payload")...)

	request, err := ReadUDPRequest(payload)
	assert.Error(err).IsNil()

	assert.Byte(request.Fragment).Equals(1)
	assert.String(request.Address).Equals("v2ray.com")
	assert.Port(request.Port).Equals(v2net.Port(80))
	assert.Bytes(request.Data.Value).Equals([]byte("Actual payload"))
}
