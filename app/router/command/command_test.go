package command_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"v2ray.com/core/app/router"
	. "v2ray.com/core/app/router/command"
	"v2ray.com/core/app/stats"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/testing/mocks"
)

func TestServiceSubscribeRoutingStats(t *testing.T) {
	c := stats.NewChannel(&stats.ChannelConfig{
		SubscriberLimit: 1,
		BufferSize:      0,
		Blocking:        true,
	})
	common.Must(c.Start())
	defer c.Close()

	lis := bufconn.Listen(1024 * 1024)
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	testCases := []*RoutingContext{
		{InboundTag: "in", OutboundTag: "out"},
		{TargetIPs: [][]byte{{1, 2, 3, 4}}, TargetPort: 8080, OutboundTag: "out"},
		{TargetDomain: "example.com", TargetPort: 443, OutboundTag: "out"},
		{SourcePort: 9999, TargetPort: 9999, OutboundTag: "out"},
		{Network: net.Network_UDP, OutboundGroupTags: []string{"outergroup", "innergroup"}, OutboundTag: "out"},
		{Protocol: "bittorrent", OutboundTag: "blocked"},
		{User: "example@v2fly.org", OutboundTag: "out"},
		{SourceIPs: [][]byte{{127, 0, 0, 1}}, Attributes: map[string]string{"attr": "value"}, OutboundTag: "out"},
	}
	errCh := make(chan error)
	nextPub := make(chan struct{})

	// Server goroutine
	go func() {
		server := grpc.NewServer()
		RegisterRoutingServiceServer(server, NewRoutingServer(nil, c))
		errCh <- server.Serve(lis)
	}()

	// Publisher goroutine
	go func() {
		publishTestCases := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			for { // Wait until there's one subscriber in routing stats channel
				if len(c.Subscribers()) > 0 {
					break
				}
				if ctx.Err() != nil {
					return ctx.Err()
				}
			}
			for _, tc := range testCases {
				c.Publish(context.Background(), AsRoutingRoute(tc))
				time.Sleep(time.Millisecond)
			}
			return nil
		}

		if err := publishTestCases(); err != nil {
			errCh <- err
		}

		// Wait for next round of publishing
		<-nextPub

		if err := publishTestCases(); err != nil {
			errCh <- err
		}
	}()

	// Client goroutine
	go func() {
		defer lis.Close()
		conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			errCh <- err
			return
		}
		defer conn.Close()
		client := NewRoutingServiceClient(conn)

		// Test retrieving all fields
		testRetrievingAllFields := func() error {
			streamCtx, streamClose := context.WithCancel(context.Background())

			// Test the unsubscription of stream works well
			defer func() {
				streamClose()
				timeOutCtx, timeout := context.WithTimeout(context.Background(), time.Second)
				defer timeout()
				for { // Wait until there's no subscriber in routing stats channel
					if len(c.Subscribers()) == 0 {
						break
					}
					if timeOutCtx.Err() != nil {
						t.Error("unexpected subscribers not decreased in channel", timeOutCtx.Err())
					}
				}
			}()

			stream, err := client.SubscribeRoutingStats(streamCtx, &SubscribeRoutingStatsRequest{})
			if err != nil {
				return err
			}

			for _, tc := range testCases {
				msg, err := stream.Recv()
				if err != nil {
					return err
				}
				if r := cmp.Diff(msg, tc, cmpopts.IgnoreUnexported(RoutingContext{})); r != "" {
					t.Error(r)
				}
			}

			// Test that double subscription will fail
			errStream, err := client.SubscribeRoutingStats(context.Background(), &SubscribeRoutingStatsRequest{
				FieldSelectors: []string{"ip", "port", "domain", "outbound"},
			})
			if err != nil {
				return err
			}
			if _, err := errStream.Recv(); err == nil {
				t.Error("unexpected successful subscription")
			}

			return nil
		}

		// Test retrieving only a subset of fields
		testRetrievingSubsetOfFields := func() error {
			streamCtx, streamClose := context.WithCancel(context.Background())
			defer streamClose()
			stream, err := client.SubscribeRoutingStats(streamCtx, &SubscribeRoutingStatsRequest{
				FieldSelectors: []string{"ip", "port", "domain", "outbound"},
			})
			if err != nil {
				return err
			}

			// Send nextPub signal to start next round of publishing
			close(nextPub)

			for _, tc := range testCases {
				msg, err := stream.Recv()
				if err != nil {
					return err
				}
				stat := &RoutingContext{ // Only a subset of stats is retrieved
					SourceIPs:         tc.SourceIPs,
					TargetIPs:         tc.TargetIPs,
					SourcePort:        tc.SourcePort,
					TargetPort:        tc.TargetPort,
					TargetDomain:      tc.TargetDomain,
					OutboundGroupTags: tc.OutboundGroupTags,
					OutboundTag:       tc.OutboundTag,
				}
				if r := cmp.Diff(msg, stat, cmpopts.IgnoreUnexported(RoutingContext{})); r != "" {
					t.Error(r)
				}
			}

			return nil
		}

		if err := testRetrievingAllFields(); err != nil {
			errCh <- err
		}
		if err := testRetrievingSubsetOfFields(); err != nil {
			errCh <- err
		}
		errCh <- nil // Client passed all tests successfully
	}()

	// Wait for goroutines to complete
	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout after 2s")
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSerivceTestRoute(t *testing.T) {
	c := stats.NewChannel(&stats.ChannelConfig{
		SubscriberLimit: 1,
		BufferSize:      16,
		Blocking:        true,
	})
	common.Must(c.Start())
	defer c.Close()

	r := new(router.Router)
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()
	common.Must(r.Init(&router.Config{
		Rule: []*router.RoutingRule{
			{
				InboundTag: []string{"in"},
				TargetTag:  &router.RoutingRule_Tag{Tag: "out"},
			},
			{
				Protocol:  []string{"bittorrent"},
				TargetTag: &router.RoutingRule_Tag{Tag: "blocked"},
			},
			{
				PortList:  &net.PortList{Range: []*net.PortRange{{From: 8080, To: 8080}}},
				TargetTag: &router.RoutingRule_Tag{Tag: "out"},
			},
			{
				SourcePortList: &net.PortList{Range: []*net.PortRange{{From: 9999, To: 9999}}},
				TargetTag:      &router.RoutingRule_Tag{Tag: "out"},
			},
			{
				Domain:    []*router.Domain{{Type: router.Domain_Domain, Value: "com"}},
				TargetTag: &router.RoutingRule_Tag{Tag: "out"},
			},
			{
				SourceGeoip: []*router.GeoIP{{CountryCode: "private", Cidr: []*router.CIDR{{Ip: []byte{127, 0, 0, 0}, Prefix: 8}}}},
				TargetTag:   &router.RoutingRule_Tag{Tag: "out"},
			},
			{
				UserEmail: []string{"example@v2fly.org"},
				TargetTag: &router.RoutingRule_Tag{Tag: "out"},
			},
			{
				Networks:  []net.Network{net.Network_UDP, net.Network_TCP},
				TargetTag: &router.RoutingRule_Tag{Tag: "out"},
			},
		},
	}, mocks.NewDNSClient(mockCtl), mocks.NewOutboundManager(mockCtl)))

	lis := bufconn.Listen(1024 * 1024)
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	errCh := make(chan error)

	// Server goroutine
	go func() {
		server := grpc.NewServer()
		RegisterRoutingServiceServer(server, NewRoutingServer(r, c))
		errCh <- server.Serve(lis)
	}()

	// Client goroutine
	go func() {
		defer lis.Close()
		conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
		if err != nil {
			errCh <- err
		}
		defer conn.Close()
		client := NewRoutingServiceClient(conn)

		testCases := []*RoutingContext{
			{InboundTag: "in", OutboundTag: "out"},
			{TargetIPs: [][]byte{{1, 2, 3, 4}}, TargetPort: 8080, OutboundTag: "out"},
			{TargetDomain: "example.com", TargetPort: 443, OutboundTag: "out"},
			{SourcePort: 9999, TargetPort: 9999, OutboundTag: "out"},
			{Network: net.Network_UDP, Protocol: "bittorrent", OutboundTag: "blocked"},
			{User: "example@v2fly.org", OutboundTag: "out"},
			{SourceIPs: [][]byte{{127, 0, 0, 1}}, Attributes: map[string]string{"attr": "value"}, OutboundTag: "out"},
		}

		// Test simple TestRoute
		testSimple := func() error {
			for _, tc := range testCases {
				route, err := client.TestRoute(context.Background(), &TestRouteRequest{RoutingContext: tc})
				if err != nil {
					return err
				}
				if r := cmp.Diff(route, tc, cmpopts.IgnoreUnexported(RoutingContext{})); r != "" {
					t.Error(r)
				}
			}
			return nil
		}

		// Test TestRoute with special options
		testOptions := func() error {
			sub, err := c.Subscribe()
			if err != nil {
				return err
			}
			for _, tc := range testCases {
				route, err := client.TestRoute(context.Background(), &TestRouteRequest{
					RoutingContext: tc,
					FieldSelectors: []string{"ip", "port", "domain", "outbound"},
					PublishResult:  true,
				})
				if err != nil {
					return err
				}
				stat := &RoutingContext{ // Only a subset of stats is retrieved
					SourceIPs:         tc.SourceIPs,
					TargetIPs:         tc.TargetIPs,
					SourcePort:        tc.SourcePort,
					TargetPort:        tc.TargetPort,
					TargetDomain:      tc.TargetDomain,
					OutboundGroupTags: tc.OutboundGroupTags,
					OutboundTag:       tc.OutboundTag,
				}
				if r := cmp.Diff(route, stat, cmpopts.IgnoreUnexported(RoutingContext{})); r != "" {
					t.Error(r)
				}
				select { // Check that routing result has been published to statistics channel
				case msg, received := <-sub:
					if route, ok := msg.(routing.Route); received && ok {
						if r := cmp.Diff(AsProtobufMessage(nil)(route), tc, cmpopts.IgnoreUnexported(RoutingContext{})); r != "" {
							t.Error(r)
						}
					} else {
						t.Error("unexpected failure in receiving published routing result for testcase", tc)
					}
				case <-time.After(100 * time.Millisecond):
					t.Error("unexpected failure in receiving published routing result", tc)
				}
			}
			return nil
		}

		if err := testSimple(); err != nil {
			errCh <- err
		}
		if err := testOptions(); err != nil {
			errCh <- err
		}
		errCh <- nil // Client passed all tests successfully
	}()

	// Wait for goroutines to complete
	select {
	case <-time.After(2 * time.Second):
		t.Fatal("Test timeout after 2s")
	case err := <-errCh:
		if err != nil {
			t.Fatal(err)
		}
	}
}
