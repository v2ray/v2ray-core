package receiver

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
