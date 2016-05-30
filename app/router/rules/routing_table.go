package rules

import (
	"sync"
	"time"
)

type RoutingEntry struct {
	tag    string
	err    error
	expire time.Time
}

func (this *RoutingEntry) Extend() {
	this.expire = time.Now().Add(time.Hour)
}

func (this *RoutingEntry) Expired() bool {
	return this.expire.Before(time.Now())
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

func (this *RoutingTable) Cleanup() {
	this.Lock()
	defer this.Unlock()

	for key, value := range this.table {
		if value.Expired() {
			delete(this.table, key)
		}
	}
}

func (this *RoutingTable) Set(destination string, tag string, err error) {
	this.Lock()
	defer this.Unlock()

	entry := &RoutingEntry{
		tag: tag,
		err: err,
	}
	entry.Extend()
	this.table[destination] = entry

	if len(this.table) > 1000 {
		go this.Cleanup()
	}
}

func (this *RoutingTable) Get(destination string) (bool, string, error) {
	this.RLock()
	defer this.RUnlock()

	entry, found := this.table[destination]
	if !found {
		return false, "", nil
	}
	entry.Extend()
	return true, entry.tag, entry.err
}
