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
	ichInUse    []proxy.InboundHandler
	ich2Recycle []proxy.InboundHandler
	lastRefresh time.Time
}

func NewInboundDetourHandlerDynamic(space app.Space, config *InboundDetourConfig) (*InboundDetourHandlerDynamic, error) {
	handler := &InboundDetourHandlerDynamic{
		space:      space,
		config:     config,
		portsInUse: make(map[v2net.Port]bool),
	}
	ichCount := config.Allocation.Concurrency
	ichArray := make([]proxy.InboundHandler, ichCount*2)
	for idx := range ichArray {
		ich, err := proxyrepo.CreateInboundHandler(config.Protocol, space, config.Settings)
		if err != nil {
			log.Error("Point: Failed to create inbound connection handler: ", err)
			return nil, err
		}
		ichArray[idx] = ich
	}
	handler.ichInUse = ichArray[:ichCount]
	handler.ich2Recycle = ichArray[ichCount:]
	return handler, nil
}

func (this *InboundDetourHandlerDynamic) pickUnusedPort() v2net.Port {
	delta := int(this.config.PortRange.To) - int(this.config.PortRange.From) + 1
	for {
		r := dice.Roll(delta)
		port := this.config.PortRange.From + v2net.Port(r)
		_, used := this.portsInUse[port]
		if !used {
			return port
		}
	}
}

func (this *InboundDetourHandlerDynamic) GetConnectionHandler() (proxy.InboundHandler, int) {
	this.RLock()
	defer this.RUnlock()
	ich := this.ichInUse[dice.Roll(len(this.ichInUse))]
	until := this.config.Allocation.Refresh - int((time.Now().Unix()-this.lastRefresh.Unix())/60/1000)
	if until < 0 {
		until = 0
	}
	return ich, int(until)
}

func (this *InboundDetourHandlerDynamic) Close() {
	this.Lock()
	defer this.Unlock()
	for _, ich := range this.ichInUse {
		ich.Close()
	}
	if this.ich2Recycle != nil {
		for _, ich := range this.ich2Recycle {
			if ich != nil {
				ich.Close()
			}
		}
	}
}

func (this *InboundDetourHandlerDynamic) refresh() error {
	this.lastRefresh = time.Now()

	for _, ich := range this.ich2Recycle {
		port2Delete := ich.Port()

		ich.Close()
		err := retry.Timed(100 /* times */, 1000 /* ms */).On(func() error {
			port := this.pickUnusedPort()
			err := ich.Listen(this.config.ListenOn, port)
			if err != nil {
				log.Error("Point: Failed to start inbound detour on port ", port, ": ", err)
				return err
			}
			this.portsInUse[port] = true
			return nil
		})
		if err != nil {
			continue
		}

		delete(this.portsInUse, port2Delete)
	}

	this.Lock()
	this.ich2Recycle, this.ichInUse = this.ichInUse, this.ich2Recycle
	this.Unlock()

	return nil
}

func (this *InboundDetourHandlerDynamic) Start() error {
	err := this.refresh()
	if err != nil {
		return err
	}

	go func() {
		for range time.Tick(time.Duration(this.config.Allocation.Refresh) * time.Minute) {
			this.refresh()
		}
	}()

	return nil
}
