package hub

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()
	atomic.StoreUint32(&o.done, 0)
}

type AwaitingConnection struct {
	conn   net.Conn
	expire time.Time
}

func (this *AwaitingConnection) Expired() bool {
	return this.expire.Before(time.Now())
}

type ConnectionCache struct {
	sync.Mutex
	cache       map[string][]*AwaitingConnection
	cleanupOnce Once
}

func NewConnectionCache() *ConnectionCache {
	return &ConnectionCache{
		cache: make(map[string][]*AwaitingConnection),
	}
}

func (this *ConnectionCache) Cleanup() {
	defer this.cleanupOnce.Reset()

	for len(this.cache) > 0 {
		time.Sleep(time.Second * 4)
		this.Lock()
		for key, value := range this.cache {
			size := len(value)
			changed := false
			for i := 0; i < size; {
				if value[i].Expired() {
					value[i].conn.Close()
					value[i] = value[size-1]
					size--
					changed = true
				} else {
					i++
				}
			}
			if changed {
				for i := size; i < len(value); i++ {
					value[i] = nil
				}
				value = value[:size]
				this.cache[key] = value
			}
		}
		this.Unlock()
	}
}

func (this *ConnectionCache) Recycle(dest string, conn net.Conn) {
	this.Lock()
	defer this.Unlock()

	aconn := &AwaitingConnection{
		conn:   conn,
		expire: time.Now().Add(time.Second * 4),
	}

	var list []*AwaitingConnection
	if v, found := this.cache[dest]; found {
		v = append(v, aconn)
		list = v
	} else {
		list = []*AwaitingConnection{aconn}
	}
	this.cache[dest] = list

	go this.cleanupOnce.Do(this.Cleanup)
}

func FindFirstValid(list []*AwaitingConnection) int {
	for idx, conn := range list {
		if !conn.Expired() {
			return idx
		}
		conn.conn.Close()
	}
	return -1
}

func (this *ConnectionCache) Get(dest string) net.Conn {
	this.Lock()
	defer this.Unlock()

	list, found := this.cache[dest]
	if !found {
		return nil
	}

	firstValid := FindFirstValid(list)
	if firstValid == -1 {
		delete(this.cache, dest)
		return nil
	}
	res := list[firstValid].conn
	list = list[firstValid+1:]
	this.cache[dest] = list
	return res
}
