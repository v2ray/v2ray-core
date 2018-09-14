package internet

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
)

var (
	effectiveTCPListener = DefaultTCPListener{}
)

type DefaultTCPListener struct{}

func (tl *DefaultTCPListener) Listen(ctx context.Context, addr *net.TCPAddr) (*net.TCPListener, error) {
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	streamSettings := StreamSettingsFromContext(ctx)
	if streamSettings != nil && streamSettings.SocketSettings != nil {
		config := streamSettings.SocketSettings
		rawConn, err := l.SyscallConn()
		if err != nil {
			return nil, err
		}
		if err := rawConn.Control(func(fd uintptr) {
			if err := applyInboundSocketOptions("tcp", fd, config); err != nil {
				newError("failed to apply socket options to incoming connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
			}
		}); err != nil {
			return nil, err
		}
	}

	return l, nil
}
