package point

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/dice"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	proxyrepo "github.com/v2ray/v2ray-core/proxy/repo"
)

type InboundDetourHandlerDynamic struct {
	sync.RWMutex
	space       app.Space
	config      *InboundDetourConfig
	portsInUse  map[v2net.Port]bool
	ichs        []proxy.InboundHandler
	lastRefresh time.Time
}

func NewInboundDetourHandlerDynamic(space app.Space, config *InboundDetourConfig) (*InboundDetourHandlerDynamic, error) {
	handler := &InboundDetourHandlerDynamic{
		space:      space,
		config:     config,
		portsInUse: make(map[v2net.Port]bool),
	}
	handler.ichs = make([]proxy.InboundHandler, config.Allocation.Concurrency)

	// To test configuration
	ich, err := proxyrepo.CreateInboundHandler(config.Protocol, space, config.Settings, &proxy.InboundHandlerMeta{
		Address: config.ListenOn,
		Port:    0,
		Tag:     config.Tag})
	if err != nil {
		log.Error("Point: Failed to create inbound connection handler: ", err)
		return nil, err
	}
	ich.Close()

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
	ich := this.ichs[dice.Roll(len(this.ichs))]
	until := this.config.Allocation.Refresh - int((time.Now().Unix()-this.lastRefresh.Unix())/60/1000)
	if until < 0 {
		until = 0
	}
	return ich, int(until)
}

func (this *InboundDetourHandlerDynamic) Close() {
	this.Lock()
	defer this.Unlock()
	for _, ich := range this.ichs {
		ich.Close()
	}
}

func (this *InboundDetourHandlerDynamic) refresh() error {
	this.lastRefresh = time.Now()

	config := this.config
	ich2Recycle := this.ichs
	newIchs := make([]proxy.InboundHandler, config.Allocation.Concurrency)

	for idx, _ := range newIchs {
		port := this.pickUnusedPort()
		ich, err := proxyrepo.CreateInboundHandler(config.Protocol, this.space, config.Settings, &proxy.InboundHandlerMeta{
			Address: config.ListenOn, Port: port, Tag: config.Tag})
		if err != nil {
			log.Error("Point: Failed to create inbound connection handler: ", err)
			return err
		}
		err = ich.Start()
		if err != nil {
			log.Error("Point: Failed to start inbound connection handler: ", err)
			return err
		}
		this.portsInUse[port] = true
		newIchs[idx] = ich
	}

	this.Lock()
	this.ichs = newIchs
	this.Unlock()

	go func(recycles []proxy.InboundHandler) {
		time.Sleep(time.Minute)
		for _, ich := range recycles {
			if ich == nil {
				continue
			}
			port := ich.Port()
			ich.Close()
			delete(this.portsInUse, port)
		}
	}(ich2Recycle)

	return nil
}

func (this *InboundDetourHandlerDynamic) Start() error {
	err := this.refresh()
	if err != nil {
		log.Error("Point: Failed to refresh dynamic allocations: ", err)
		return err
	}

	go func() {
		for {
			time.Sleep(time.Duration(this.config.Allocation.Refresh) * time.Minute)
			err := this.refresh()
			if err != nil {
				log.Error("Point: Failed to refresh dynamic allocations: ", err)
			}
		}
	}()

	return nil
}
