package commander

//go:generate errorgen

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/signal/done"
)

// Commander is a V2Ray feature that provides gRPC methods to external clients.
type Commander struct {
	sync.Mutex
	server *grpc.Server
	config Config
	v      *core.Instance
	ohm    core.OutboundHandlerManager
}

// NewCommander creates a new Commander based on the given config.
func NewCommander(ctx context.Context, config *Config) (*Commander, error) {
	v := core.MustFromContext(ctx)
	c := &Commander{
		config: *config,
		ohm:    v.OutboundHandlerManager(),
		v:      v,
	}
	if err := v.RegisterFeature((*Commander)(nil), c); err != nil {
		return nil, err
	}
	return c, nil
}

// Type implements common.HasType.
func (c *Commander) Type() interface{} {
	return (*Commander)(nil)
}

// Start implements common.Runnable.
func (c *Commander) Start() error {
	c.Lock()
	c.server = grpc.NewServer()
	for _, rawConfig := range c.config.Service {
		config, err := rawConfig.GetInstance()
		if err != nil {
			return err
		}
		rawService, err := core.CreateObject(c.v, config)
		if err != nil {
			return err
		}
		service, ok := rawService.(Service)
		if !ok {
			return newError("not a Service.")
		}
		service.Register(c.server)
	}
	c.Unlock()

	listener := &OutboundListener{
		buffer: make(chan net.Conn, 4),
		done:   done.New(),
	}

	go func() {
		if err := c.server.Serve(listener); err != nil {
			newError("failed to start grpc server").Base(err).AtError().WriteToLog()
		}
	}()

	if err := c.ohm.RemoveHandler(context.Background(), c.config.Tag); err != nil {
		newError("failed to remove existing handler").WriteToLog()
	}

	return c.ohm.AddHandler(context.Background(), &Outbound{
		tag:      c.config.Tag,
		listener: listener,
	})
}

// Close implements common.Closable.
func (c *Commander) Close() error {
	c.Lock()
	defer c.Unlock()

	if c.server != nil {
		c.server.Stop()
		c.server = nil
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		return NewCommander(ctx, cfg.(*Config))
	}))
}
