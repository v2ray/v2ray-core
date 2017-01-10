package receiver

import (
	"net"

	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type refresher struct {
	strategy   *AllocationStrategy
	portsInUse []v2net.Port
}

func (r *refresher) Refresh(s *StreamReceiver) {
}

func (r *refresher) Interval() time.Duration {
	switch r.strategy.Type {
	case AllocationStrategy_Random:
		return time.Minute * time.Duration(r.strategy.GetRefreshValue())
	default:
		return 0
	}
}

type StreamReceiver struct {
	config *StreamReceiverConfig
	proxy  *proxy.InboundHandler

	listeners []net.Listener
	refresher refresher
}

func (s *StreamReceiver) Start() {
	s.refresher.Refresh(s)
	interval := s.refresher.Interval()
	if interval == 0 {
		return
	}

	go func() {
		for {
			time.Sleep(s.refresher.Interval())
			s.refresher.Refresh(s)
		}
	}()
}
