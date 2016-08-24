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

func (this *ServerList) AddServer(server *ServerSpec) {
	this.Lock()
	defer this.Unlock()

	this.servers = append(this.servers, server)
}

func (this *ServerList) Size() uint32 {
	this.RLock()
	defer this.RUnlock()

	return uint32(len(this.servers))
}

func (this *ServerList) GetServer(idx uint32) *ServerSpec {
	this.RLock()
	defer this.RUnlock()

	for {
		if idx >= uint32(len(this.servers)) {
			return nil
		}

		server := this.servers[idx]
		if !server.IsValid() {
			this.RemoveServer(idx)
			continue
		}

		return server
	}
}

// Private: Visible for testing.
func (this *ServerList) RemoveServer(idx uint32) {
	n := len(this.servers)
	this.servers[idx] = this.servers[n-1]
	this.servers = this.servers[:n-1]
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

func (this *RoundRobinServerPicker) PickServer() *ServerSpec {
	this.Lock()
	defer this.Unlock()

	next := this.nextIndex
	server := this.serverlist.GetServer(next)
	if server == nil {
		next = 0
		server = this.serverlist.GetServer(0)
	}
	next++
	if next >= this.serverlist.Size() {
		next = 0
	}
	this.nextIndex = next

	return server
}
