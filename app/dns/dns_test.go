package dns_test

import (
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/app/dns"
	apptesting "github.com/v2ray/v2ray-core/app/testing"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	v2testing "github.com/v2ray/v2ray-core/testing"
)

func TestDnsAdd(t *testing.T) {
	v2testing.Current(t)

	domain := "v2ray.com"
	cache := dns.NewCache(nil)
	ip := cache.Get(&apptesting.Context{}, domain)
	netassert.IP(ip).IsNil()

	cache.Add(&apptesting.Context{}, domain, []byte{1, 2, 3, 4})
	ip = cache.Get(&apptesting.Context{}, domain)
	netassert.IP(ip).Equals(net.IP([]byte{1, 2, 3, 4}))
}
