package udp

import (
	"sync"
	"time"

	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type UDPResponseCallback func(destination v2net.Destination, payload *buf.Buffer)

type TimedInboundRay struct {
	name       string
	inboundRay ray.InboundRay
	accessed   chan bool
	server     *UDPServer
	sync.RWMutex
}

func NewTimedInboundRay(name string, inboundRay ray.InboundRay, server *UDPServer) *TimedInboundRay {
	r := &TimedInboundRay{
		name:       name,
		inboundRay: inboundRay,
		accessed:   make(chan bool, 1),
		server:     server,
	}
	go r.Monitor()
	return r
}

func (v *TimedInboundRay) Monitor() {
	for {
		time.Sleep(time.Second * 16)
		select {
		case <-v.accessed:
		default:
			// Ray not accessed for a while, assuming communication is dead.
			v.RLock()
			if v.server == nil {
				v.RUnlock()
				return
			}
			v.server.RemoveRay(v.name)
			v.RUnlock()
			v.Release()
			return
		}
	}
}

func (v *TimedInboundRay) InboundInput() ray.OutputStream {
	v.RLock()
	defer v.RUnlock()
	if v.inboundRay == nil {
		return nil
	}
	select {
	case v.accessed <- true:
	default:
	}
	return v.inboundRay.InboundInput()
}

func (v *TimedInboundRay) InboundOutput() ray.InputStream {
	v.RLock()
	defer v.RUnlock()
	if v.inboundRay == nil {
		return nil
	}
	select {
	case v.accessed <- true:
	default:
	}
	return v.inboundRay.InboundOutput()
}

func (v *TimedInboundRay) Release() {
	log.Debug("UDP Server: Releasing TimedInboundRay: ", v.name)
	v.Lock()
	defer v.Unlock()
	if v.server == nil {
		return
	}
	v.server = nil
	v.inboundRay.InboundInput().Close()
	v.inboundRay.InboundOutput().Release()
	v.inboundRay = nil
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

func (v *UDPServer) RemoveRay(name string) {
	v.Lock()
	defer v.Unlock()
	delete(v.conns, name)
}

func (v *UDPServer) locateExistingAndDispatch(name string, payload *buf.Buffer) bool {
	log.Debug("UDP Server: Locating existing connection for ", name)
	v.RLock()
	defer v.RUnlock()
	if entry, found := v.conns[name]; found {
		outputStream := entry.InboundInput()
		if outputStream == nil {
			return false
		}
		err := outputStream.Write(payload)
		if err != nil {
			go entry.Release()
			return false
		}
		return true
	}
	return false
}

func (v *UDPServer) Dispatch(session *proxy.SessionInfo, payload *buf.Buffer, callback UDPResponseCallback) {
	source := session.Source
	destination := session.Destination

	// TODO: Add user to destString
	destString := source.String() + "-" + destination.String()
	log.Debug("UDP Server: Dispatch request: ", destString)
	if v.locateExistingAndDispatch(destString, payload) {
		return
	}

	log.Info("UDP Server: establishing new connection for ", destString)
	inboundRay := v.packetDispatcher.DispatchToOutbound(session)
	timedInboundRay := NewTimedInboundRay(destString, inboundRay, v)
	outputStream := timedInboundRay.InboundInput()
	if outputStream != nil {
		outputStream.Write(payload)
	}

	v.Lock()
	v.conns[destString] = timedInboundRay
	v.Unlock()
	go v.handleConnection(timedInboundRay, source, callback)
}

func (v *UDPServer) handleConnection(inboundRay *TimedInboundRay, source v2net.Destination, callback UDPResponseCallback) {
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
