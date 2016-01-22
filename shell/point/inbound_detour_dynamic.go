package point

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/dice"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy"
	proxyrepo "github.com/v2ray/v2ray-core/proxy/repo"
)

type InboundDetourHandlerDynamic struct {
	sync.RWMutex
	space       app.Space
	config      *InboundDetourConfig
	portsInUse  map[v2net.Port]bool
	ichInUse    []*InboundConnectionHandlerWithPort
	ich2Recycle []*InboundConnectionHandlerWithPort
	lastRefresh time.Time
	started     bool
}

func NewInboundDetourHandlerDynamic(space app.Space, config *InboundDetourConfig) (*InboundDetourHandlerDynamic, error) {
	handler := &InboundDetourHandlerDynamic{
		space:      space,
		config:     config,
		portsInUse: make(map[v2net.Port]bool),
	}
	if err := handler.refresh(); err != nil {
		return nil, err
	}
	return handler, nil
}

func (this *InboundDetourHandlerDynamic) refresh() error {
	this.Lock()
	defer this.Unlock()

	this.ich2Recycle = this.ichInUse
	if this.ich2Recycle != nil {
		time.AfterFunc(10*time.Second, func() {
			for _, ich := range this.ich2Recycle {
				if ich != nil {
					ich.handler.Close()
					delete(this.portsInUse, ich.port)
				}
			}
		})
	}

	ichCount := this.config.Allocation.Concurrency
	// TODO: check ichCount
	this.ichInUse = make([]*InboundConnectionHandlerWithPort, ichCount)
	for idx, _ := range this.ichInUse {
		port := this.pickUnusedPort()
		ich, err := proxyrepo.CreateInboundConnectionHandler(this.config.Protocol, this.space, this.config.Settings)
		if err != nil {
			log.Error("Point: Failed to create inbound connection handler: ", err)
			return err
		}
		this.ichInUse[idx] = &InboundConnectionHandlerWithPort{
			port:    port,
			handler: ich,
		}
	}
	if this.started {
		this.Start()
	}

	this.lastRefresh = time.Now()
	time.AfterFunc(time.Duration(this.config.Allocation.Refresh)*time.Minute, func() {
		this.refresh()
	})

	return nil
}

func (this *InboundDetourHandlerDynamic) pickUnusedPort() v2net.Port {
	delta := int(this.config.PortRange.To) - int(this.config.PortRange.From) + 1
	for {
		r := dice.Roll(delta)
		port := this.config.PortRange.From + v2net.Port(r)
		_, used := this.portsInUse[port]
		if !used {
			this.portsInUse[port] = true
			return port
		}
	}
}

func (this *InboundDetourHandlerDynamic) GetConnectionHandler() (proxy.InboundConnectionHandler, int) {
	this.RLock()
	defer this.RUnlock()
	ich := this.ichInUse[dice.Roll(len(this.ichInUse))]
	until := this.config.Allocation.Refresh - int((time.Now().Unix()-this.lastRefresh.Unix())/60/1000)
	if until < 0 {
		until = 0
	}
	return ich.handler, int(until)
}

func (this *InboundDetourHandlerDynamic) Close() {
	this.Lock()
	defer this.Unlock()
	for _, ich := range this.ichInUse {
		ich.handler.Close()
	}
	if this.ich2Recycle != nil {
		for _, ich := range this.ich2Recycle {
			if ich != nil && ich.handler != nil {
				ich.handler.Close()
			}
		}
	}
}

func (this *InboundDetourHandlerDynamic) Start() error {
	for _, ich := range this.ichInUse {
		err := retry.Timed(100 /* times */, 100 /* ms */).On(func() error {
			err := ich.handler.Listen(ich.port)
			if err != nil {
				log.Error("Point: Failed to start inbound detour on port ", ich.port, ": ", err)
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	this.started = true
	return nil
}
