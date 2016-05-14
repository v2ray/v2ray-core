package hub

import (
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type UDPResponseCallback func(destination v2net.Destination, payload *alloc.Buffer)

type TimedInboundRay struct {
	name       string
	inboundRay ray.InboundRay
	accessed   chan bool
	server     *UDPServer
	sync.RWMutex
}

func NewTimedInboundRay(name string, inboundRay ray.InboundRay) *TimedInboundRay {
	r := &TimedInboundRay{
		name:       name,
		inboundRay: inboundRay,
		accessed:   make(chan bool),
	}
	go r.Monitor()
	return r
}

func (this *TimedInboundRay) Monitor() {
	for {
		time.Sleep(time.Second * 16)
		select {
		case <-this.accessed:
		default:
			// Ray not accessed for a while, assuming communication is dead.
			this.Release()
			return
		}
	}
}

func (this *TimedInboundRay) InboundInput() ray.OutputStream {
	this.RLock()
	defer this.RUnlock()
	if this.inboundRay == nil {
		return nil
	}
	select {
	case this.accessed <- true:
	default:
	}
	return this.inboundRay.InboundInput()
}

func (this *TimedInboundRay) InboundOutput() ray.InputStream {
	this.RLock()
	defer this.RUnlock()
	if this.inboundRay == nil {
		return nil
	}
	select {
	case this.accessed <- true:
	default:
	}
	return this.inboundRay.InboundOutput()
}

func (this *TimedInboundRay) Release() {
	log.Debug("UDP Server: Releasing TimedInboundRay: ", this.name)
	this.Lock()
	defer this.Unlock()
	if this.server == nil {
		return
	}
	this.server.RemoveRay(this.name)
	this.server = nil
	this.inboundRay.InboundInput().Close()
	this.inboundRay.InboundOutput().Release()
	this.inboundRay = nil
}

type UDPServer struct {
	sync.RWMutex
	conns            map[string]*TimedInboundRay
	packetDispatcher dispatcher.PacketDispatcher
}

func NewUDPServer(packetDispatcher dispatcher.PacketDispatcher) *UDPServer {
	return &UDPServer{
		conns:            make(map[string]*TimedInboundRay),
		packetDispatcher: packetDispatcher,
	}
}

func (this *UDPServer) RemoveRay(name string) {
	this.Lock()
	defer this.Unlock()
	delete(this.conns, name)
}

func (this *UDPServer) locateExistingAndDispatch(name string, payload *alloc.Buffer) bool {
	log.Debug("UDP Server: Locating existing connection for ", name)
	this.RLock()
	defer this.RUnlock()
	if entry, found := this.conns[name]; found {
		outputStream := entry.InboundInput()
		if outputStream == nil {
			return false
		}
		err := outputStream.Write(payload)
		if err != nil {
			go this.RemoveRay(name)
			return false
		}
		return true
	}
	return false
}

func (this *UDPServer) Dispatch(source v2net.Destination, destination v2net.Destination, payload *alloc.Buffer, callback UDPResponseCallback) {
	destString := source.NetAddr() + "-" + destination.NetAddr()
	log.Debug("UDP Server: Dispatch request: ", destString)
	if this.locateExistingAndDispatch(destString, payload) {
		return
	}

	log.Info("UDP Server: establishing new connection for ", destString)
	inboundRay := this.packetDispatcher.DispatchToOutbound(destination)
	timedInboundRay := NewTimedInboundRay(destString, inboundRay)
	outputStream := timedInboundRay.InboundInput()
	if outputStream != nil {
		outputStream.Write(payload)
	}

	this.Lock()
	this.conns[destString] = timedInboundRay
	this.Unlock()
	go this.handleConnection(timedInboundRay, source, callback)
}

func (this *UDPServer) handleConnection(inboundRay *TimedInboundRay, source v2net.Destination, callback UDPResponseCallback) {
	for {
		inputStream := inboundRay.InboundOutput()
		if inputStream == nil {
			break
		}
		data, err := inboundRay.InboundOutput().Read()
		if err != nil {
			break
		}
		callback(source, data)
	}
	inboundRay.Release()
}
