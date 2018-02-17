package policy

import (
	"time"

	"v2ray.com/core"
)

// Duration converts Second to time.Duration.
func (s *Second) Duration() time.Duration {
	if s == nil {
		return 0
	}
	return time.Second * time.Duration(s.Value)
}

func (p *Policy) ToCorePolicy() core.Policy {
	var cp core.Policy
	if p.Timeout != nil {
		cp.Timeouts.ConnectionIdle = p.Timeout.ConnectionIdle.Duration()
		cp.Timeouts.Handshake = p.Timeout.Handshake.Duration()
		cp.Timeouts.DownlinkOnly = p.Timeout.DownlinkOnly.Duration()
		cp.Timeouts.UplinkOnly = p.Timeout.UplinkOnly.Duration()
	}
	return cp
}
