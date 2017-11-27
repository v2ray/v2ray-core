package policy

import (
	"time"
)

func (s *Second) Duration() time.Duration {
	return time.Second * time.Duration(s.Value)
}

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
