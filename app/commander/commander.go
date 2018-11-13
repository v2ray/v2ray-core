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
	"v2ray.com/core/features/outbound"
)

// Commander is a V2Ray feature that provides gRPC methods to external clients.
type Commander struct {
	sync.Mutex
	server   *grpc.Server
	services []Service
	ohm      outbound.Manager
	tag      string
}

// NewCommander creates a new Commander based on the given config.
func NewCommander(ctx context.Context, config *Config) (*Commander, error) {
	c := &Commander{
		tag: config.Tag,
	}

	common.Must(core.RequireFeatures(ctx, func(om outbound.Manager) {
		c.ohm = om
	}))

	for _, rawConfig := range config.Service {
		config, err := rawConfig.GetInstance()
		if err != nil {
			return nil, err
		}
		rawService, err := common.CreateObject(ctx, config)
		if err != nil {
			return nil, err
		}
		service, ok := rawService.(Service)
		if !ok {
			return nil, newError("not a Service.")
		}
		c.services = append(c.services, service)
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
	for _, service := range c.services {
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

	if err := c.ohm.RemoveHandler(context.Background(), c.tag); err != nil {
		newError("failed to remove existing handler").WriteToLog()
	}

	return c.ohm.AddHandler(context.Background(), &Outbound{
		tag:      c.tag,
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
