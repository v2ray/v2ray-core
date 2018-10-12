package stats

import "v2ray.com/core/features"

type Counter interface {
	Value() int64
	Set(int64) int64
	Add(int64) int64
}

type Manager interface {
	features.Feature

	RegisterCounter(string) (Counter, error)
	GetCounter(string) Counter
}

// GetOrRegisterCounter tries to get the StatCounter first. If not exist, it then tries to create a new counter.
func GetOrRegisterCounter(m Manager, name string) (Counter, error) {
	counter := m.GetCounter(name)
	if counter != nil {
		return counter, nil
	}

	return m.RegisterCounter(name)
}
