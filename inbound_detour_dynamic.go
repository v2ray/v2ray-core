package core

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/proxy"
)

type InboundDetourHandlerDynamic struct {
	sync.RWMutex
	space       app.Space
	config      *InboundConnectionConfig
	portsInUse  map[net.Port]bool
	ichs        []proxy.InboundHandler
	ich2Recyle  []proxy.InboundHandler
	lastRefresh time.Time
	ctx         context.Context
}

func NewInboundDetourHandlerDynamic(ctx context.Context, config *InboundConnectionConfig) (*InboundDetourHandlerDynamic, error) {
	space := app.SpaceFromContext(ctx)
	handler := &InboundDetourHandlerDynamic{
		space:      space,
		config:     config,
		portsInUse: make(map[net.Port]bool),
		ctx:        ctx,
	}
	handler.ichs = make([]proxy.InboundHandler, config.GetAllocationStrategyValue().GetConcurrencyValue())

	// To test configuration
	ichConfig, err := config.GetTypedSettings()
	if err != nil {
		return nil, err
	}
	ich, err := proxy.CreateInboundHandler(proxy.ContextWithInboundMeta(ctx, &proxy.InboundHandlerMeta{
		Address:                config.GetListenOnValue(),
		Port:                   0,
		Tag:                    config.Tag,
		StreamSettings:         config.StreamSettings,
		AllowPassiveConnection: config.AllowPassiveConnection,
	}), ichConfig)
	if err != nil {
		log.Error("Point: Failed to create inbound connection handler: ", err)
		return nil, err
	}
	ich.Close()

	return handler, nil
}

func (v *InboundDetourHandlerDynamic) pickUnusedPort() net.Port {
	delta := int(v.config.PortRange.To) - int(v.config.PortRange.From) + 1
	for {
		r := dice.Roll(delta)
		port := v.config.PortRange.FromPort() + net.Port(r)
		_, used := v.portsInUse[port]
		if !used {
			return port
		}
	}
}

func (v *InboundDetourHandlerDynamic) GetConnectionHandler() (proxy.InboundHandler, int) {
	v.RLock()
	defer v.RUnlock()
	ich := v.ichs[dice.Roll(len(v.ichs))]
	until := int(v.config.GetAllocationStrategyValue().GetRefreshValue()) - int((time.Now().Unix()-v.lastRefresh.Unix())/60/1000)
	if until < 0 {
		until = 0
	}
	return ich, int(until)
}

func (v *InboundDetourHandlerDynamic) Close() {
	v.Lock()
	defer v.Unlock()
	for _, ich := range v.ichs {
		ich.Close()
	}
}

func (v *InboundDetourHandlerDynamic) RecyleHandles() {
	if v.ich2Recyle != nil {
		for _, ich := range v.ich2Recyle {
			if ich == nil {
				continue
			}
			port := ich.Port()
			ich.Close()
			delete(v.portsInUse, port)
		}
		v.ich2Recyle = nil
	}
}

func (v *InboundDetourHandlerDynamic) refresh() error {
	v.lastRefresh = time.Now()

	config := v.config
	v.ich2Recyle = v.ichs
	newIchs := make([]proxy.InboundHandler, config.GetAllocationStrategyValue().GetConcurrencyValue())

	for idx := range newIchs {
		err := retry.Timed(5, 100).On(func() error {
			port := v.pickUnusedPort()
			ichConfig, _ := config.GetTypedSettings()
			ich, err := proxy.CreateInboundHandler(proxy.ContextWithInboundMeta(v.ctx, &proxy.InboundHandlerMeta{
				Address: config.GetListenOnValue(),
				Port:    port, Tag: config.Tag,
				StreamSettings: config.StreamSettings}), ichConfig)
			if err != nil {
				delete(v.portsInUse, port)
				return err
			}
			err = ich.Start()
			if err != nil {
				delete(v.portsInUse, port)
				return err
			}
			v.portsInUse[port] = true
			newIchs[idx] = ich
			return nil
		})
		if err != nil {
			log.Error("Point: Failed to create inbound connection handler: ", err)
			return err
		}
	}

	v.Lock()
	v.ichs = newIchs
	v.Unlock()

	return nil
}

func (v *InboundDetourHandlerDynamic) Start() error {
	err := v.refresh()
	if err != nil {
		log.Error("Point: Failed to refresh dynamic allocations: ", err)
		return err
	}

	go func() {
		for {
			time.Sleep(time.Duration(v.config.GetAllocationStrategyValue().GetRefreshValue())*time.Minute - 1)
			v.RecyleHandles()
			err := v.refresh()
			if err != nil {
				log.Error("Point: Failed to refresh dynamic allocations: ", err)
			}
			time.Sleep(time.Minute)
		}
	}()

	return nil
}
