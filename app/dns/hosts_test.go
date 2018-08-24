package dns_test

import (
	"testing"

	. "v2ray.com/core/app/dns"
	. "v2ray.com/ext/assert"
)

func TestStaticHosts(t *testing.T) {
	assert := With(t)

	pb := []*Config_HostMapping{
		{
			Type:   DomainMatchingType_Full,
			Domain: "v2ray.com",
			Ip: [][]byte{
				{1, 1, 1, 1},
			},
		},
		{
			Type:   DomainMatchingType_Subdomain,
			Domain: "v2ray.cn",
			Ip: [][]byte{
				{2, 2, 2, 2},
			},
		},
	}

	hosts, err := NewStaticHosts(pb, nil)
	assert(err, IsNil)

	{
		ips := hosts.LookupIP("v2ray.com")
		assert(len(ips), Equals, 1)
		assert([]byte(ips[0]), Equals, []byte{1, 1, 1, 1})
	}

	{
		ips := hosts.LookupIP("www.v2ray.cn")
		assert(len(ips), Equals, 1)
		assert([]byte(ips[0]), Equals, []byte{2, 2, 2, 2})
	}
}
