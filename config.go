package core

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
)

func (v *AllocationStrategy) GetConcurrencyValue() uint32 {
	if v == nil || v.Concurrency == nil {
		return 3
	}
	return v.Concurrency.Value
}

func (v *AllocationStrategy) GetRefreshValue() uint32 {
	if v == nil || v.Refresh == nil {
		return 5
	}
	return v.Refresh.Value
}

func (v *InboundConnectionConfig) GetAllocationStrategyValue() *AllocationStrategy {
	if v.AllocationStrategy == nil {
		return &AllocationStrategy{}
	}
	return v.AllocationStrategy
}

func (v *InboundConnectionConfig) GetListenOnValue() net.Address {
	if v.GetListenOn() == nil {
		return net.AnyIP
	}
	return v.ListenOn.AsAddress()
}

func (v *InboundConnectionConfig) GetTypedSettings() (interface{}, error) {
	if v.GetSettings() == nil {
		return nil, common.ErrBadConfiguration
	}
	return v.GetSettings().GetInstance()
}

func (v *OutboundConnectionConfig) GetTypedSettings() (interface{}, error) {
	if v.GetSettings() == nil {
		return nil, common.ErrBadConfiguration
	}
	return v.GetSettings().GetInstance()
}

func (v *OutboundConnectionConfig) GetSendThroughValue() net.Address {
	if v.GetSendThrough() == nil {
		return net.AnyIP
	}
	return v.SendThrough.AsAddress()
}
