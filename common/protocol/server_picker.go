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

func (sl *ServerList) AddServer(server *ServerSpec) {
	sl.Lock()
	defer sl.Unlock()

	sl.servers = append(sl.servers, server)
}

func (sl *ServerList) Size() uint32 {
	sl.RLock()
	defer sl.RUnlock()

	return uint32(len(sl.servers))
}

func (sl *ServerList) GetServer(idx uint32) *ServerSpec {
	sl.Lock()
	defer sl.Unlock()

	for {
		if idx >= uint32(len(sl.servers)) {
			return nil
		}

		server := sl.servers[idx]
		if !server.IsValid() {
			sl.removeServer(idx)
			continue
		}

		return server
	}
}

func (sl *ServerList) removeServer(idx uint32) {
	n := len(sl.servers)
	sl.servers[idx] = sl.servers[n-1]
	sl.servers = sl.servers[:n-1]
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

func (p *RoundRobinServerPicker) PickServer() *ServerSpec {
	p.Lock()
	defer p.Unlock()

	next := p.nextIndex
	server := p.serverlist.GetServer(next)
	if server == nil {
		next = 0
		server = p.serverlist.GetServer(0)
	}
	next++
	if next >= p.serverlist.Size() {
		next = 0
	}
	p.nextIndex = next

	return server
}
