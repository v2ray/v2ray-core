package internet

import (
	"context"
	"syscall"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
)

var (
	effectiveListener = DefaultListener{}
)

type DefaultListener struct{}

func (*DefaultListener) Listen(ctx context.Context, addr net.Addr) (net.Listener, error) {
	var lc net.ListenConfig

	sockopt := getSocketSettings(ctx)
	if sockopt != nil {
		lc.Control = func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if err := applyInboundSocketOptions(network, fd, sockopt); err != nil {
					newError("failed to apply socket options to incoming connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
				}
			})
		}
	}

	return lc.Listen(ctx, addr.Network(), addr.String())
}

func (*DefaultListener) ListenPacket(ctx context.Context, addr net.Addr) (net.PacketConn, error) {
	var lc net.ListenConfig

	sockopt := getSocketSettings(ctx)
	if sockopt != nil {
		lc.Control = func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				if err := applyInboundSocketOptions(network, fd, sockopt); err != nil {
					newError("failed to apply socket options to incoming connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
				}
			})
		}
	}

	return lc.ListenPacket(ctx, addr.Network(), addr.String())
}
