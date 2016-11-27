package router

import (
	"sync"
	"time"
)

type RoutingEntry struct {
	tag    string
	err    error
	expire time.Time
}

func (v *RoutingEntry) Extend() {
	v.expire = time.Now().Add(time.Hour)
}

func (v *RoutingEntry) Expired() bool {
	return v.expire.Before(time.Now())
}

type RoutingTable struct {
	sync.RWMutex
	table map[string]*RoutingEntry
}

func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		table: make(map[string]*RoutingEntry),
	}
}

func (v *RoutingTable) Cleanup() {
	v.Lock()
	defer v.Unlock()

	for key, value := range v.table {
		if value.Expired() {
			delete(v.table, key)
		}
	}
}

func (v *RoutingTable) Set(destination string, tag string, err error) {
	v.Lock()
	defer v.Unlock()

	entry := &RoutingEntry{
		tag: tag,
		err: err,
	}
	entry.Extend()
	v.table[destination] = entry

	if len(v.table) > 1000 {
		go v.Cleanup()
	}
}

func (v *RoutingTable) Get(destination string) (bool, string, error) {
	v.RLock()
	defer v.RUnlock()

	entry, found := v.table[destination]
	if !found {
		return false, "", nil
	}
	entry.Extend()
	return true, entry.tag, entry.err
}
