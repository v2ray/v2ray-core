package hub

import (
	"sync"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type UDPResponseCallback func(packet v2net.Packet)

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

func (this *UDPServer) locateExistingAndDispatch(dest string, packet v2net.Packet) bool {
	this.RLock()
	defer this.RUnlock()
	if entry, found := this.conns[dest]; found {
		entry.inboundRay.InboundInput().Write(packet.Chunk())
		return true
	}
	return false
}

func (this *UDPServer) Dispatch(source v2net.Destination, packet v2net.Packet, callback UDPResponseCallback) {
	destString := source.String() + "-" + packet.Destination().NetAddr()
	if this.locateExistingAndDispatch(destString, packet) {
		return
	}

	this.Lock()
	inboundRay := this.packetDispatcher.DispatchToOutbound(v2net.NewPacket(packet.Destination(), packet.Chunk(), true))
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
		callback(v2net.NewPacket(source, data, false))
	}
	this.Lock()
	delete(this.conns, destString)
	this.Unlock()
}
