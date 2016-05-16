package dns

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/hub"

	"github.com/miekg/dns"
)

const (
	DefaultTTL = uint32(3600)
)

type ARecord struct {
	IPs    []net.IP
	Expire time.Time
}

type NameServer interface {
	QueryA(domain string) <-chan *ARecord
}

type PendingRequest struct {
	expire   time.Time
	response chan<- *ARecord
}

type UDPNameServer struct {
	sync.Mutex
	address   v2net.Destination
	requests  map[uint16]*PendingRequest
	udpServer *hub.UDPServer
}

func NewUDPNameServer(address v2net.Destination, dispatcher dispatcher.PacketDispatcher) *UDPNameServer {
	s := &UDPNameServer{
		address:   address,
		requests:  make(map[uint16]*PendingRequest),
		udpServer: hub.NewUDPServer(dispatcher),
	}
	go s.Cleanup()
	return s
}

// @Private
func (this *UDPNameServer) Cleanup() {
	for {
		time.Sleep(time.Second * 60)
		expiredRequests := make([]uint16, 0, 16)
		now := time.Now()
		this.Lock()
		for id, r := range this.requests {
			if r.expire.Before(now) {
				expiredRequests = append(expiredRequests, id)
				close(r.response)
			}
		}
		for _, id := range expiredRequests {
			delete(this.requests, id)
		}
		this.Unlock()
		expiredRequests = nil
	}
}

// @Private
func (this *UDPNameServer) AssignUnusedID(response chan<- *ARecord) uint16 {
	var id uint16
	this.Lock()
	for {
		id = uint16(rand.Intn(65536))
		if _, found := this.requests[id]; found {
			continue
		}
		log.Debug("DNS: Add pending request id ", id)
		this.requests[id] = &PendingRequest{
			expire:   time.Now().Add(time.Second * 16),
			response: response,
		}
		break
	}
	this.Unlock()
	return id
}

// @Private
func (this *UDPNameServer) HandleResponse(dest v2net.Destination, payload *alloc.Buffer) {
	msg := new(dns.Msg)
	err := msg.Unpack(payload.Value)
	if err != nil {
		log.Warning("DNS: Failed to parse DNS response: ", err)
		return
	}
	record := &ARecord{
		IPs: make([]net.IP, 0, 16),
	}
	id := msg.Id
	ttl := DefaultTTL

	this.Lock()
	request, found := this.requests[id]
	if !found {
		this.Unlock()
		return
	}
	delete(this.requests, id)
	this.Unlock()

	for _, rr := range msg.Answer {
		if a, ok := rr.(*dns.A); ok {
			record.IPs = append(record.IPs, a.A)
			if a.Hdr.Ttl < ttl {
				ttl = a.Hdr.Ttl
			}
		}
	}
	record.Expire = time.Now().Add(time.Second * time.Duration(ttl))

	request.response <- record
	close(request.response)
}

func (this *UDPNameServer) QueryA(domain string) <-chan *ARecord {
	response := make(chan *ARecord)

	buffer := alloc.NewBuffer()
	msg := new(dns.Msg)
	msg.Id = this.AssignUnusedID(response)
	msg.RecursionDesired = true
	msg.Question = []dns.Question{
		dns.Question{
			Name:   dns.Fqdn(domain),
			Qtype:  dns.TypeA,
			Qclass: dns.ClassINET,
		},
		dns.Question{
			Name:   dns.Fqdn(domain),
			Qtype:  dns.TypeAAAA,
			Qclass: dns.ClassINET,
		},
	}

	writtenBuffer, _ := msg.PackBuffer(buffer.Value)
	buffer.Slice(0, len(writtenBuffer))

	fakeDestination := v2net.UDPDestination(v2net.LocalHostIP, v2net.Port(53))
	this.udpServer.Dispatch(fakeDestination, this.address, buffer, this.HandleResponse)

	return response
}
