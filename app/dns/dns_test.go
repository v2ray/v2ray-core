package dns_test

import (
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/app/dns"
	dnstesting "github.com/v2ray/v2ray-core/app/dns/testing"
	apptesting "github.com/v2ray/v2ray-core/app/testing"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	v2testing "github.com/v2ray/v2ray-core/testing"
)

func TestDnsAdd(t *testing.T) {
	v2testing.Current(t)

	domain := "v2ray.com"
	cache := dns.NewCache(&dnstesting.CacheConfig{
		TrustedTags: map[string]bool{
			"testtag": true,
		},
	})
	ip := cache.Get(&apptesting.Context{}, domain)
	netassert.IP(ip).IsNil()

	cache.Add(&apptesting.Context{CallerTagValue: "notvalidtag"}, domain, []byte{1, 2, 3, 4})
	ip = cache.Get(&apptesting.Context{}, domain)
	netassert.IP(ip).IsNil()

	cache.Add(&apptesting.Context{CallerTagValue: "testtag"}, domain, []byte{1, 2, 3, 4})
	ip = cache.Get(&apptesting.Context{}, domain)
	netassert.IP(ip).Equals(net.IP([]byte{1, 2, 3, 4}))
}
