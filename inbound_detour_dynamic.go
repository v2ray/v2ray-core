package core

import (
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/proxy"
	proxyregistry "v2ray.com/core/proxy/registry"
)

type InboundDetourHandlerDynamic struct {
	sync.RWMutex
	space       app.Space
	config      *InboundConnectionConfig
	portsInUse  map[v2net.Port]bool
	ichs        []proxy.InboundHandler
	ich2Recyle  []proxy.InboundHandler
	lastRefresh time.Time
}

func NewInboundDetourHandlerDynamic(space app.Space, config *InboundConnectionConfig) (*InboundDetourHandlerDynamic, error) {
	handler := &InboundDetourHandlerDynamic{
		space:      space,
		config:     config,
		portsInUse: make(map[v2net.Port]bool),
	}
	handler.ichs = make([]proxy.InboundHandler, config.GetAllocationStrategyValue().Concurrency.GetValue())

	// To test configuration
	ich, err := proxyregistry.CreateInboundHandler(config.Protocol, space, config.Settings, &proxy.InboundHandlerMeta{
		Address:                config.ListenOn.AsAddress(),
		Port:                   0,
		Tag:                    config.Tag,
		StreamSettings:         config.StreamSettings,
		AllowPassiveConnection: config.AllowPassiveConnection,
	})
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
		port := this.config.PortRange.FromPort() + v2net.Port(r)
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
	until := int(this.config.GetAllocationStrategyValue().Refresh.GetValue()) - int((time.Now().Unix()-this.lastRefresh.Unix())/60/1000)
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

func (this *InboundDetourHandlerDynamic) RecyleHandles() {
	if this.ich2Recyle != nil {
		for _, ich := range this.ich2Recyle {
			if ich == nil {
				continue
			}
			port := ich.Port()
			ich.Close()
			delete(this.portsInUse, port)
		}
		this.ich2Recyle = nil
	}
}

func (this *InboundDetourHandlerDynamic) refresh() error {
	this.lastRefresh = time.Now()

	config := this.config
	this.ich2Recyle = this.ichs
	newIchs := make([]proxy.InboundHandler, config.GetAllocationStrategyValue().Concurrency.GetValue())

	for idx := range newIchs {
		err := retry.Timed(5, 100).On(func() error {
			port := this.pickUnusedPort()
			ich, err := proxyregistry.CreateInboundHandler(config.Protocol, this.space, config.Settings, &proxy.InboundHandlerMeta{
				Address: config.ListenOn.AsAddress(), Port: port, Tag: config.Tag, StreamSettings: config.StreamSettings})
			if err != nil {
				delete(this.portsInUse, port)
				return err
			}
			err = ich.Start()
			if err != nil {
				delete(this.portsInUse, port)
				return err
			}
			this.portsInUse[port] = true
			newIchs[idx] = ich
			return nil
		})
		if err != nil {
			log.Error("Point: Failed to create inbound connection handler: ", err)
			return err
		}
	}

	this.Lock()
	this.ichs = newIchs
	this.Unlock()

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
			time.Sleep(time.Duration(this.config.GetAllocationStrategyValue().Refresh.GetValue())*time.Minute - 1)
			this.RecyleHandles()
			err := this.refresh()
			if err != nil {
				log.Error("Point: Failed to refresh dynamic allocations: ", err)
			}
			time.Sleep(time.Minute)
		}
	}()

	return nil
}
