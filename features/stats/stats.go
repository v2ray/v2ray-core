package stats

//go:generate errorgen

import "v2ray.com/core/features"

// Counter is the interface for stats counters.
//
// v2ray:api:stable
type Counter interface {
	// Value is the current value of the counter.
	Value() int64
	// Set sets a new value to the counter, and returns the previous one.
	Set(int64) int64
	// Add adds a value to the current counter value, and returns the previous value.
	Add(int64) int64
}

// Manager is the interface for stats manager.
//
// v2ray:api:stable
type Manager interface {
	features.Feature

	// RegisterCounter registers a new counter to the manager. The identifier string must not be emtpy, and unique among other counters.
	RegisterCounter(string) (Counter, error)
	UnregisterCounter(string) error
	// GetCounter returns a counter by its identifier.
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

// ManagerType returns the type of Manager interface. Can be used to implement common.HasType.
//
// v2ray:api:stable
func ManagerType() interface{} {
	return (*Manager)(nil)
}

// NoopManager is an implementation of Manager, which doesn't has actual functionalities.
type NoopManager struct{}

// Type implements common.HasType.
func (NoopManager) Type() interface{} {
	return ManagerType()
}

// RegisterCounter implements Manager.
func (NoopManager) RegisterCounter(string) (Counter, error) {
	return nil, newError("not implemented")
}

// UnregisterCounter implements Manager.
func (NoopManager) UnregisterCounter(string) error {
	return newError("not implemented")
}

// GetCounter implements Manager.
func (NoopManager) GetCounter(string) Counter {
	return nil
}

// Start implements common.Runnable.
func (NoopManager) Start() error { return nil }

// Close implements common.Closable.
func (NoopManager) Close() error { return nil }
