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

// OverrideWith overrides current Policy with another one.
func (p *Policy) OverrideWith(another *Policy) {
	if another.Timeout != nil {
		if another.Timeout.Handshake != nil {
			p.Timeout.Handshake = another.Timeout.Handshake
		}
		if another.Timeout.ConnectionIdle != nil {
			p.Timeout.ConnectionIdle = another.Timeout.ConnectionIdle
		}
		if another.Timeout.UplinkOnly != nil {
			p.Timeout.UplinkOnly = another.Timeout.UplinkOnly
		}
		if another.Timeout.DownlinkOnly != nil {
			p.Timeout.DownlinkOnly = another.Timeout.DownlinkOnly
		}
	}
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
