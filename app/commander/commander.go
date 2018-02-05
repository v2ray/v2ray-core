package commander

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg commander -path App,Commander

import (
	"context"
	"net"
	"sync"

	"google.golang.org/grpc"
	"v2ray.com/core"
)

type Commander struct {
	sync.Mutex
	server    *grpc.Server
	config    Config
	ohm       core.OutboundHandlerManager
	callbacks []core.ServiceRegistryCallback
}

func (c *Commander) RegisterService(callback core.ServiceRegistryCallback) {
	c.Lock()
	defer c.Unlock()

	if callback == nil {
		return
	}

	c.callbacks = append(c.callbacks, callback)
}

func (c *Commander) Start() error {
	c.Lock()
	c.server = grpc.NewServer()
	for _, callback := range c.callbacks {
		callback(c.server)
	}
	c.Unlock()

	listener := &OutboundListener{
		buffer: make(chan net.Conn, 4),
	}

	c.server.Serve(listener)

	c.ohm.RemoveHandler(context.Background(), c.config.Tag)
	c.ohm.AddHandler(context.Background(), &CommanderOutbound{
		tag:      c.config.Tag,
		listener: listener,
	})
	return nil
}

func (c *Commander) Close() {
	c.Lock()
	defer c.Unlock()

	if c.server != nil {
		c.server.Stop()
		c.server = nil
	}
}
