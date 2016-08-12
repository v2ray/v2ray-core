package dns

import (
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/dice"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/internet/udp"

	"github.com/miekg/dns"
)

const (
	DefaultTTL       = uint32(3600)
	CleanupInterval  = time.Second * 120
	CleanupThreshold = 512
)

var (
	pseudoDestination = v2net.UDPDestination(v2net.LocalHostIP, v2net.Port(53))
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
	address     v2net.Destination
	requests    map[uint16]*PendingRequest
	udpServer   *udp.UDPServer
	nextCleanup time.Time
}

func NewUDPNameServer(address v2net.Destination, dispatcher dispatcher.PacketDispatcher) *UDPNameServer {
	s := &UDPNameServer{
		address:  address,
		requests: make(map[uint16]*PendingRequest),
		udpServer: udp.NewUDPServer(&proxy.InboundHandlerMeta{
			AllowPassiveConnection: false,
		}, dispatcher),
	}
	return s
}

// @Private
func (this *UDPNameServer) Cleanup() {
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

// @Private
func (this *UDPNameServer) AssignUnusedID(response chan<- *ARecord) uint16 {
	var id uint16
	this.Lock()
	if len(this.requests) > CleanupThreshold && this.nextCleanup.Before(time.Now()) {
		this.nextCleanup = time.Now().Add(CleanupInterval)
		go this.Cleanup()
	}

	for {
		id = uint16(dice.Roll(65536))
		if _, found := this.requests[id]; found {
			continue
		}
		log.Debug("DNS: Add pending request id ", id)
		this.requests[id] = &PendingRequest{
			expire:   time.Now().Add(time.Second * 8),
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
	log.Debug("DNS: Handling response for id ", id, " content: ", msg.String())

	this.Lock()
	request, found := this.requests[id]
	if !found {
		this.Unlock()
		return
	}
	delete(this.requests, id)
	this.Unlock()

	for _, rr := range msg.Answer {
		switch rr := rr.(type) {
		case *dns.A:
			record.IPs = append(record.IPs, rr.A)
			if rr.Hdr.Ttl < ttl {
				ttl = rr.Hdr.Ttl
			}
		case *dns.AAAA:
			record.IPs = append(record.IPs, rr.AAAA)
			if rr.Hdr.Ttl < ttl {
				ttl = rr.Hdr.Ttl
			}
		}
	}
	record.Expire = time.Now().Add(time.Second * time.Duration(ttl))

	request.response <- record
	close(request.response)
}

func (this *UDPNameServer) BuildQueryA(domain string, id uint16) *alloc.Buffer {
	buffer := alloc.NewBuffer()
	msg := new(dns.Msg)
	msg.Id = id
	msg.RecursionDesired = true
	msg.Question = []dns.Question{
		{
			Name:   dns.Fqdn(domain),
			Qtype:  dns.TypeA,
			Qclass: dns.ClassINET,
		}}

	writtenBuffer, _ := msg.PackBuffer(buffer.Value)
	buffer.Slice(0, len(writtenBuffer))

	return buffer
}

func (this *UDPNameServer) DispatchQuery(payload *alloc.Buffer) {
	this.udpServer.Dispatch(pseudoDestination, this.address, payload, this.HandleResponse)
}

func (this *UDPNameServer) QueryA(domain string) <-chan *ARecord {
	response := make(chan *ARecord, 1)
	id := this.AssignUnusedID(response)

	this.DispatchQuery(this.BuildQueryA(domain, id))

	go func() {
		for i := 0; i < 2; i++ {
			time.Sleep(time.Second)
			this.Lock()
			_, found := this.requests[id]
			this.Unlock()
			if found {
				this.DispatchQuery(this.BuildQueryA(domain, id))
			} else {
				break
			}
		}
	}()

	return response
}

type LocalNameServer struct {
}

func (this *LocalNameServer) QueryA(domain string) <-chan *ARecord {
	response := make(chan *ARecord, 1)

	go func() {
		defer close(response)

		ips, err := net.LookupIP(domain)
		if err != nil {
			log.Info("DNS: Failed to lookup IPs for domain ", domain)
			return
		}

		response <- &ARecord{
			IPs:    ips,
			Expire: time.Now().Add(time.Second * time.Duration(DefaultTTL)),
		}
	}()

	return response
}
