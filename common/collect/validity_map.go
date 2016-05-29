package collect

import (
	"sync"
	"sync/atomic"

	"github.com/v2ray/v2ray-core/common"
)

type Validity interface {
	common.Releasable
	IsValid() bool
}

type ValidityMap struct {
	sync.RWMutex
	cache   map[string]Validity
	opCount int32
}

func NewValidityMap(cleanupIntervalSec int) *ValidityMap {
	instance := &ValidityMap{
		cache: make(map[string]Validity),
	}
	return instance
}

func (this *ValidityMap) cleanup() {
	type entry struct {
		key   string
		value Validity
	}

	entry2Remove := make([]entry, 0, 128)
	this.RLock()
	for key, value := range this.cache {
		if !value.IsValid() {
			entry2Remove = append(entry2Remove, entry{
				key:   key,
				value: value,
			})
		}
	}
	this.RUnlock()

	for _, e := range entry2Remove {
		if !e.value.IsValid() {
			this.Lock()
			delete(this.cache, e.key)
			this.Unlock()
		}
	}
}

func (this *ValidityMap) Set(key string, value Validity) {
	this.Lock()
	this.cache[key] = value
	this.Unlock()
	opCount := atomic.AddInt32(&this.opCount, 1)
	if opCount > 1000 {
		atomic.StoreInt32(&this.opCount, 0)
		go this.cleanup()
	}
}

func (this *ValidityMap) Get(key string) Validity {
	this.RLock()
	defer this.RUnlock()
	if value, found := this.cache[key]; found {
		return value
	}
	return nil
}
