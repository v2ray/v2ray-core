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

type controller func(network, address string, fd uintptr) error

type DefaultListener struct {
	contollers []controller
}

func getControlFunc(ctx context.Context, sockopt *SocketConfig, contollers []controller) func(network, address string, c syscall.RawConn) error {
	return func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			if sockopt != nil {
				if err := applyInboundSocketOptions(network, fd, sockopt); err != nil {
					newError("failed to apply socket options to incoming connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
				}
			}

			for _, controller := range contollers {
				if err := controller(network, address, fd); err != nil {
					newError("failed to apply external controller").Base(err).WriteToLog(session.ExportIDToError(ctx))
				}
			}
		})
	}
}

func (dl *DefaultListener) Listen(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.Listener, error) {
	var lc net.ListenConfig

	if sockopt != nil || len(dl.contollers) > 0 {
		lc.Control = getControlFunc(ctx, sockopt, dl.contollers)
	}

	return lc.Listen(ctx, addr.Network(), addr.String())
}

func (dl *DefaultListener) ListenPacket(ctx context.Context, addr net.Addr, sockopt *SocketConfig) (net.PacketConn, error) {
	var lc net.ListenConfig

	if sockopt != nil || len(dl.contollers) > 0 {
		lc.Control = getControlFunc(ctx, sockopt, dl.contollers)
	}

	return lc.ListenPacket(ctx, addr.Network(), addr.String())
}

// RegisterListenerController adds a controller to the effective system listener.
// The controller can be used to operate on file descriptors before they are put into use.
func RegisterListenerController(controller func(network, address string, fd uintptr) error) error {
	if controller == nil {
		return newError("nil listener controller")
	}

	effectiveListener.contollers = append(effectiveListener.contollers, controller)
	return nil
}
