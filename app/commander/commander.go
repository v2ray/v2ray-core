package commander

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg commander -path App,Commander

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/signal"
)

type Commander struct {
	sync.Mutex
	server *grpc.Server
	config Config
	v      *core.Instance
	ohm    core.OutboundHandlerManager
}

func NewCommander(ctx context.Context, config *Config) (*Commander, error) {
	v := core.FromContext(ctx)
	if v == nil {
		return nil, newError("V is not in context.")
	}
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

func (c *Commander) Type() interface{} {
	return (*Commander)(nil)
}

func (c *Commander) Start() error {
	c.Lock()
	c.server = grpc.NewServer()
	for _, rawConfig := range c.config.Service {
		config, err := rawConfig.GetInstance()
		if err != nil {
			return err
		}
		rawService, err := c.v.CreateObject(config)
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
		done:   signal.NewDone(),
	}

	go func() {
		if err := c.server.Serve(listener); err != nil {
			newError("failed to start grpc server").Base(err).AtError().WriteToLog()
		}
	}()

	c.ohm.RemoveHandler(context.Background(), c.config.Tag)
	c.ohm.AddHandler(context.Background(), &CommanderOutbound{
		tag:      c.config.Tag,
		listener: listener,
	})
	return nil
}

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
