package command_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	. "github.com/v2ray/v2ray-core/proxy/vmess/command"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestCacheDnsIPv4(t *testing.T) {
	v2testing.Current(t)

	cd := &CacheDns{
		Address: v2net.IPAddress([]byte{1, 2, 3, 4}),
	}

	buffer := alloc.NewBuffer().Clear()
	defer buffer.Release()

	nBytes, err := cd.Marshal(buffer)
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(buffer.Len())

	cd2 := &CacheDns{}
	err = cd2.Unmarshal(buffer.Value)
	assert.Error(err).IsNil()
	netassert.Address(cd.Address).Equals(cd2.Address)
}
