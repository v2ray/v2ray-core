package core

import (
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
)

func (v *AllocationStrategyConcurrency) GetValue() uint32 {
	if v == nil {
		return 3
	}
	return v.Value
}

func (v *AllocationStrategyRefresh) GetValue() uint32 {
	if v == nil {
		return 5
	}
	return v.Value
}

func (v *InboundConnectionConfig) GetAllocationStrategyValue() *AllocationStrategy {
	if v.AllocationStrategy == nil {
		return &AllocationStrategy{}
	}
	return v.AllocationStrategy
}

func (v *InboundConnectionConfig) GetListenOnValue() v2net.Address {
	if v.GetListenOn() == nil {
		return v2net.AnyIP
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

func (v *OutboundConnectionConfig) GetSendThroughValue() v2net.Address {
	if v.GetSendThrough() == nil {
		return v2net.AnyIP
	}
	return v.SendThrough.AsAddress()
}
