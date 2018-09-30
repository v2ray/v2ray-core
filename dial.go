package core

import (
	"context"

	"v2ray.com/core/common/net"
)

// Dial provides an easy way for upstream caller to create net.Conn through V2Ray.
// It dispatches the request to the given destination by the given V2Ray instance.
// Since it is under a proxy context, the LocalAddr() and RemoteAddr() in returned net.Conn
// will not show real addresses being used for communication.
func Dial(ctx context.Context, v *Instance, dest net.Destination) (net.Conn, error) {
	r, err := v.Dispatcher().Dispatch(ctx, dest)
	if err != nil {
		return nil, err
	}
	return net.NewConnection(net.ConnectionInputMulti(r.Writer), net.ConnectionOutputMulti(r.Reader)), nil
}
