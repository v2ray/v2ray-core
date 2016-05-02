package hub

import (
	"sync"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type UDPResponseCallback func(destination v2net.Destination, payload *alloc.Buffer)

type connEntry struct {
	inboundRay ray.InboundRay
	callback   UDPResponseCallback
}

type UDPServer struct {
	sync.RWMutex
	conns            map[string]*connEntry
	packetDispatcher dispatcher.PacketDispatcher
}

func NewUDPServer(packetDispatcher dispatcher.PacketDispatcher) *UDPServer {
	return &UDPServer{
		conns:            make(map[string]*connEntry),
		packetDispatcher: packetDispatcher,
	}
}

func (this *UDPServer) locateExistingAndDispatch(dest string, payload *alloc.Buffer) bool {
	this.RLock()
	defer this.RUnlock()
	if entry, found := this.conns[dest]; found {
		entry.inboundRay.InboundInput().Write(payload)
		return true
	}
	return false
}

func (this *UDPServer) Dispatch(source v2net.Destination, destination v2net.Destination, payload *alloc.Buffer, callback UDPResponseCallback) {
	destString := source.String() + "-" + destination.NetAddr()
	if this.locateExistingAndDispatch(destString, payload) {
		return
	}

	this.Lock()
	inboundRay := this.packetDispatcher.DispatchToOutbound(destination)
	inboundRay.InboundInput().Write(payload)

	this.conns[destString] = &connEntry{
		inboundRay: inboundRay,
		callback:   callback,
	}
	this.Unlock()
	go this.handleConnection(destString, inboundRay, source, callback)
}

func (this *UDPServer) handleConnection(destString string, inboundRay ray.InboundRay, source v2net.Destination, callback UDPResponseCallback) {
	for {
		data, err := inboundRay.InboundOutput().Read()
		if err != nil {
			break
		}
		callback(source, data)
	}
	this.Lock()
	inboundRay.InboundInput().Release()
	inboundRay.InboundOutput().Release()
	delete(this.conns, destString)
	this.Unlock()
}
