package internet

import (
	"context"
	"net"
	"time"

	v2net "v2ray.com/core/common/net"
)

var (
	effectiveSystemDialer SystemDialer = DefaultSystemDialer{}
)

type SystemDialer interface {
	Dial(ctx context.Context, source v2net.Address, destination v2net.Destination) (net.Conn, error)
}

type DefaultSystemDialer struct {
}

func (DefaultSystemDialer) Dial(ctx context.Context, src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   time.Second * 60,
		DualStack: true,
	}
	if src != nil && src != v2net.AnyIP {
		var addr net.Addr
		if dest.Network == v2net.Network_TCP {
			addr = &net.TCPAddr{
				IP:   src.IP(),
				Port: 0,
			}
		} else {
			addr = &net.UDPAddr{
				IP:   src.IP(),
				Port: 0,
			}
		}
		dialer.LocalAddr = addr
	}
	return dialer.DialContext(ctx, dest.Network.SystemString(), dest.NetAddr())
}

type SystemDialerAdapter interface {
	Dial(network string, address string) (net.Conn, error)
}

type SimpleSystemDialer struct {
	adapter SystemDialerAdapter
}

func WithAdapter(dialer SystemDialerAdapter) SystemDialer {
	return &SimpleSystemDialer{
		adapter: dialer,
	}
}

func (v *SimpleSystemDialer) Dial(ctx context.Context, src v2net.Address, dest v2net.Destination) (net.Conn, error) {
	return v.adapter.Dial(dest.Network.SystemString(), dest.NetAddr())
}

// UseAlternativeSystemDialer replaces the current system dialer with a given one.
// Caller must ensure there is no race condition.
func UseAlternativeSystemDialer(dialer SystemDialer) {
	if dialer == nil {
		effectiveSystemDialer = DefaultSystemDialer{}
	}
	effectiveSystemDialer = dialer
}
