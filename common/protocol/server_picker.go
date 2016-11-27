package protocol

import (
	"sync"
)

type ServerList struct {
	sync.RWMutex
	servers []*ServerSpec
}

func NewServerList() *ServerList {
	return &ServerList{}
}

func (v *ServerList) AddServer(server *ServerSpec) {
	v.Lock()
	defer v.Unlock()

	v.servers = append(v.servers, server)
}

func (v *ServerList) Size() uint32 {
	v.RLock()
	defer v.RUnlock()

	return uint32(len(v.servers))
}

func (v *ServerList) GetServer(idx uint32) *ServerSpec {
	v.RLock()
	defer v.RUnlock()

	for {
		if idx >= uint32(len(v.servers)) {
			return nil
		}

		server := v.servers[idx]
		if !server.IsValid() {
			v.RemoveServer(idx)
			continue
		}

		return server
	}
}

// Private: Visible for testing.
func (v *ServerList) RemoveServer(idx uint32) {
	n := len(v.servers)
	v.servers[idx] = v.servers[n-1]
	v.servers = v.servers[:n-1]
}

type ServerPicker interface {
	PickServer() *ServerSpec
}

type RoundRobinServerPicker struct {
	sync.Mutex
	serverlist *ServerList
	nextIndex  uint32
}

func NewRoundRobinServerPicker(serverlist *ServerList) *RoundRobinServerPicker {
	return &RoundRobinServerPicker{
		serverlist: serverlist,
		nextIndex:  0,
	}
}

func (v *RoundRobinServerPicker) PickServer() *ServerSpec {
	v.Lock()
	defer v.Unlock()

	next := v.nextIndex
	server := v.serverlist.GetServer(next)
	if server == nil {
		next = 0
		server = v.serverlist.GetServer(0)
	}
	next++
	if next >= v.serverlist.Size() {
		next = 0
	}
	v.nextIndex = next

	return server
}
