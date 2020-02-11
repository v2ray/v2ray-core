package dns_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	. "v2ray.com/core/app/dns"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
)

func TestStaticHosts(t *testing.T) {
	pb := []*Config_HostMapping{
		{
			Pattern: "fv2ray.com",
			Ip: [][]byte{
				{1, 1, 1, 1},
			},
		},
		{
			Pattern: "dv2ray.cn",
			Ip: [][]byte{
				{2, 2, 2, 2},
			},
		},
		{
			Pattern: "dbaidu.com",
			Ip: [][]byte{
				{127, 0, 0, 1},
			},
		},
	}

	hosts, err := NewStaticHosts(pb, nil, make(map[string][]string))
	common.Must(err)

	{
		ips := hosts.LookupIP("v2ray.com", IPOption{
			IPv4Enable: true,
			IPv6Enable: true,
		})
		if len(ips) != 1 {
			t.Error("expect 1 IP, but got ", len(ips))
		}
		if diff := cmp.Diff([]byte(ips[0].IP()), []byte{1, 1, 1, 1}); diff != "" {
			t.Error(diff)
		}
	}

	{
		ips := hosts.LookupIP("www.v2ray.cn", IPOption{
			IPv4Enable: true,
			IPv6Enable: true,
		})
		if len(ips) != 1 {
			t.Error("expect 1 IP, but got ", len(ips))
		}
		if diff := cmp.Diff([]byte(ips[0].IP()), []byte{2, 2, 2, 2}); diff != "" {
			t.Error(diff)
		}
	}

	{
		ips := hosts.LookupIP("baidu.com", IPOption{
			IPv4Enable: false,
			IPv6Enable: true,
		})
		if len(ips) != 1 {
			t.Error("expect 1 IP, but got ", len(ips))
		}
		if diff := cmp.Diff([]byte(ips[0].IP()), []byte(net.LocalHostIPv6.IP())); diff != "" {
			t.Error(diff)
		}
	}
}
