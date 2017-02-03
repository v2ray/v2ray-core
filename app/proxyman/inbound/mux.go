package inbound

import (
	"context"
	"errors"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/ray"
)

type mux struct {
	dispatcher dispatcher.Interface
}

func newMux(ctx context.Context) *mux {
	m := &mux{}
	space := app.SpaceFromContext(ctx)
	space.OnInitialize(func() error {
		d := dispatcher.FromSpace(space)
		if d == nil {
			return errors.New("Proxyman|DefaultInboundHandler: No dispatcher in space.")
		}
		m.dispatcher = d
		return nil
	})
	return m
}

func (m *mux) Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error) {
	return m.dispatcher.Dispatch(ctx, dest)
}
