package dns_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "v2ray.com/core/app/dns"
	"v2ray.com/core/common"
)

func TestStaticHosts(t *testing.T) {
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
	common.Must(err)

	{
		ips := hosts.LookupIP("v2ray.com")
		if len(ips) != 1 {
			t.Error("expect 1 IP, but got ", len(ips))
		}
		if diff := cmp.Diff([]byte(ips[0]), []byte{1, 1, 1, 1}); diff != "" {
			t.Error(diff)
		}
	}

	{
		ips := hosts.LookupIP("www.v2ray.cn")
		if len(ips) != 1 {
			t.Error("expect 1 IP, but got ", len(ips))
		}
		if diff := cmp.Diff([]byte(ips[0]), []byte{2, 2, 2, 2}); diff != "" {
			t.Error(diff)
		}
	}
}
