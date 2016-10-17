package core

import (
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
)

func (this *AllocationStrategyConcurrency) GetValue() uint32 {
	if this == nil {
		return 3
	}
	return this.Value
}

func (this *AllocationStrategyRefresh) GetValue() uint32 {
	if this == nil {
		return 5
	}
	return this.Value
}

func (this *InboundConnectionConfig) GetAllocationStrategyValue() *AllocationStrategy {
	if this.AllocationStrategy == nil {
		return &AllocationStrategy{}
	}
	return this.AllocationStrategy
}

func (this *InboundConnectionConfig) GetListenOnValue() v2net.Address {
	if this.GetListenOn() == nil {
		return v2net.AnyIP
	}
	return this.ListenOn.AsAddress()
}

func (this *InboundConnectionConfig) GetTypedSettings() (interface{}, error) {
	if this.GetSettings() == nil {
		return nil, common.ErrBadConfiguration
	}
	return this.GetSettings().GetInstance()
}

func (this *OutboundConnectionConfig) GetTypedSettings() (interface{}, error) {
	if this.GetSettings() == nil {
		return nil, common.ErrBadConfiguration
	}
	return this.GetSettings().GetInstance()
}

func (this *OutboundConnectionConfig) GetSendThroughValue() v2net.Address {
	if this.GetSendThrough() == nil {
		return v2net.AnyIP
	}
	return this.SendThrough.AsAddress()
}
