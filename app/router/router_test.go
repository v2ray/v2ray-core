package router_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/testing/mocks"
)

func TestSimpleRouter(t *testing.T) {
	config := &Config{
		Rule: []*RoutingRule{
			{
				TargetTag: &RoutingRule_Tag{
					Tag: "test",
				},
				NetworkList: &net.NetworkList{
					Network: []net.Network{net.Network_TCP},
				},
			},
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockDns := mocks.NewDNSClient(mockCtl)

	r := new(Router)
	common.Must(r.Init(config, mockDns, nil))

	ctx := withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("v2ray.com"), 80)})
	tag, err := r.PickRoute(ctx)
	common.Must(err)
	if tag != "test" {
		t.Error("expect tag 'test', bug actually ", tag)
	}
}

func TestIPOnDemand(t *testing.T) {
	config := &Config{
		DomainStrategy: Config_IpOnDemand,
		Rule: []*RoutingRule{
			{
				TargetTag: &RoutingRule_Tag{
					Tag: "test",
				},
				Cidr: []*CIDR{
					{
						Ip:     []byte{192, 168, 0, 0},
						Prefix: 16,
					},
				},
			},
		},
	}

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockDns := mocks.NewDNSClient(mockCtl)
	mockDns.EXPECT().LookupIP(gomock.Eq("v2ray.com")).Return([]net.IP{{192, 168, 0, 1}}, nil).AnyTimes()

	r := new(Router)
	common.Must(r.Init(config, mockDns, nil))

	ctx := withOutbound(&session.Outbound{Target: net.TCPDestination(net.DomainAddress("v2ray.com"), 80)})
	tag, err := r.PickRoute(ctx)
	common.Must(err)
	if tag != "test" {
		t.Error("expect tag 'test', bug actually ", tag)
	}
}
