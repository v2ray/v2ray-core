package core

import (
	"v2ray.com/core/proxy/registry"
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

func (this *InboundConnectionConfig) GetTypedSettings() (interface{}, error) {
	return registry.MarshalInboundConfig(this.Protocol, this.Settings)
}
