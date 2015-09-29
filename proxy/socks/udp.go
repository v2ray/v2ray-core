package socks

import (
	"math"
	"math/rand"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/common/collect"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/socks/protocol"
)

const (
	bufferSize = 2 * 1024
)

type portMap struct {
	access       sync.Mutex
	data         map[uint16]*net.UDPAddr
	removedPorts *collect.TimedQueue
}

func newPortMap() *portMap {
	m := &portMap{
		access:       sync.Mutex{},
		data:         make(map[uint16]*net.UDPAddr),
		removedPorts: collect.NewTimedQueue(1),
	}
	go m.removePorts(m.removedPorts.RemovedEntries())
	return m
}

func (m *portMap) assignAddressToken(addr *net.UDPAddr) uint16 {
	for {
		token := uint16(rand.Intn(math.MaxUint16))
		if _, used := m.data[token]; !used {
			m.access.Lock()
			if _, used = m.data[token]; !used {
				m.data[token] = addr
				m.access.Unlock()
				return token
			}
			m.access.Unlock()
		}
	}
}

func (m *portMap) removePorts(removedPorts <-chan interface{}) {
	for {
		rawToken := <-removedPorts
		m.access.Lock()
		delete(m.data, rawToken.(uint16))
		m.access.Unlock()
	}
}

func (m *portMap) popPort(token uint16) *net.UDPAddr {
	m.access.Lock()
	defer m.access.Unlock()
	addr, exists := m.data[token]
	if !exists {
		return nil
	}
	delete(m.data, token)
	return addr
}

var (
	ports = newPortMap()

	udpConn *net.UDPConn
)

func (server *SocksServer) ListenUDP(port uint16) error {
	addr := &net.UDPAddr{
		IP:   net.IP{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Error("Socks failed to listen UDP on port %d: %v", port, err)
		return err
	}

	go server.AcceptPackets(conn)
	udpConn = conn
	return nil
}

func (server *SocksServer) AcceptPackets(conn *net.UDPConn) error {
	for {
		buffer := make([]byte, 0, bufferSize)
		nBytes, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Error("Socks failed to read UDP packets: %v", err)
			return err
		}
		request, err := protocol.ReadUDPRequest(buffer[:nBytes])
		if err != nil {
			log.Error("Socks failed to parse UDP request: %v", err)
			return err
		}
		if request.Fragment != 0 {
			// TODO handle fragments
			continue
		}

		token := ports.assignAddressToken(addr)

		udpPacket := v2net.NewUDPPacket(request.Destination(), request.Data, token)
		server.vPoint.DispatchToOutbound(udpPacket)
	}
}

func (server *SocksServer) Dispatch(packet v2net.Packet) {
	if udpPacket, ok := packet.(*v2net.UDPPacket); ok {
		token := udpPacket.Token()
		addr := ports.popPort(token)
		if udpConn != nil {
			udpConn.WriteToUDP(udpPacket.Chunk(), addr)
		}
	}
	// We don't expect TCP Packets here
}
