package websocket

import (
	"net"
	"sync"
	"time"

	"v2ray.com/core/common/log"
	"v2ray.com/core/common/signal"
)

type AwaitingConnection struct {
	conn   *wsconn
	expire time.Time
}

func (v *AwaitingConnection) Expired() bool {
	return v.expire.Before(time.Now())
}

type ConnectionCache struct {
	sync.Mutex
	cache       map[string][]*AwaitingConnection
	cleanupOnce signal.Once
}

func NewConnectionCache() *ConnectionCache {
	return &ConnectionCache{
		cache: make(map[string][]*AwaitingConnection),
	}
}

func (v *ConnectionCache) Cleanup() {
	defer v.cleanupOnce.Reset()

	for len(v.cache) > 0 {
		time.Sleep(time.Second * 7)
		v.Lock()
		for key, value := range v.cache {
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
				v.cache[key] = value
			}
		}
		v.Unlock()
	}
}

func (v *ConnectionCache) Recycle(dest string, conn *wsconn) {
	v.Lock()
	defer v.Unlock()

	aconn := &AwaitingConnection{
		conn:   conn,
		expire: time.Now().Add(time.Second * 7),
	}

	var list []*AwaitingConnection
	if val, found := v.cache[dest]; found {
		val = append(val, aconn)
		list = val
	} else {
		list = []*AwaitingConnection{aconn}
	}
	v.cache[dest] = list

	go v.cleanupOnce.Do(v.Cleanup)
}

func FindFirstValid(list []*AwaitingConnection) int {
	for idx, conn := range list {
		if !conn.Expired() && !conn.conn.connClosing {
			return idx
		}
		go conn.conn.Close()
	}
	return -1
}

func (v *ConnectionCache) Get(dest string) net.Conn {
	v.Lock()
	defer v.Unlock()

	list, found := v.cache[dest]
	if !found {
		return nil
	}

	firstValid := FindFirstValid(list)
	if firstValid == -1 {
		delete(v.cache, dest)
		return nil
	}
	res := list[firstValid].conn
	list = list[firstValid+1:]
	v.cache[dest] = list
	log.Debug("WS:Conn Cache used.")
	return res
}
