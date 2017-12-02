package server

import (
	"context"
	"sync"
	"time"

	"github.com/miekg/dns"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet/udp"
)

const (
	CleanupInterval  = time.Second * 120
	CleanupThreshold = 512
)

var (
	multiQuestionDNS = map[net.Address]bool{
		net.IPAddress([]byte{8, 8, 8, 8}): true,
		net.IPAddress([]byte{8, 8, 4, 4}): true,
		net.IPAddress([]byte{9, 9, 9, 9}): true,
	}
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
	address     net.Destination
	requests    map[uint16]*PendingRequest
	udpServer   *udp.Dispatcher
	nextCleanup time.Time
}

func NewUDPNameServer(address net.Destination, dispatcher dispatcher.Interface) *UDPNameServer {
	s := &UDPNameServer{
		address:   address,
		requests:  make(map[uint16]*PendingRequest),
		udpServer: udp.NewDispatcher(dispatcher),
	}
	return s
}

func (s *UDPNameServer) Cleanup() {
	expiredRequests := make([]uint16, 0, 16)
	now := time.Now()
	s.Lock()
	for id, r := range s.requests {
		if r.expire.Before(now) {
			expiredRequests = append(expiredRequests, id)
			close(r.response)
		}
	}
	for _, id := range expiredRequests {
		delete(s.requests, id)
	}
	s.Unlock()
}

func (s *UDPNameServer) AssignUnusedID(response chan<- *ARecord) uint16 {
	var id uint16
	s.Lock()
	if len(s.requests) > CleanupThreshold && s.nextCleanup.Before(time.Now()) {
		s.nextCleanup = time.Now().Add(CleanupInterval)
		go s.Cleanup()
	}

	for {
		id = dice.RollUint16()
		if _, found := s.requests[id]; found {
			continue
		}
		log.Trace(newError("add pending request id ", id).AtDebug())
		s.requests[id] = &PendingRequest{
			expire:   time.Now().Add(time.Second * 8),
			response: response,
		}
		break
	}
	s.Unlock()
	return id
}

func (s *UDPNameServer) HandleResponse(payload *buf.Buffer) {
	msg := new(dns.Msg)
	err := msg.Unpack(payload.Bytes())
	if err == dns.ErrTruncated {
		log.Trace(newError("truncated message received. DNS server should still work. If you see anything abnormal, please submit an issue to v2ray-core.").AtWarning())
	} else if err != nil {
		log.Trace(newError("failed to parse DNS response").Base(err).AtWarning())
		return
	}
	record := &ARecord{
		IPs: make([]net.IP, 0, 16),
	}
	id := msg.Id
	ttl := uint32(3600) // an hour
	log.Trace(newError("handling response for id ", id, " content: ", msg).AtDebug())

	s.Lock()
	request, found := s.requests[id]
	if !found {
		s.Unlock()
		return
	}
	delete(s.requests, id)
	s.Unlock()

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

func (s *UDPNameServer) BuildQueryA(domain string, id uint16) *buf.Buffer {
	msg := new(dns.Msg)
	msg.Id = id
	msg.RecursionDesired = true
	msg.Question = []dns.Question{
		{
			Name:   dns.Fqdn(domain),
			Qtype:  dns.TypeA,
			Qclass: dns.ClassINET,
		}}
	if multiQuestionDNS[s.address.Address] {
		msg.Question = append(msg.Question, dns.Question{
			Name:   dns.Fqdn(domain),
			Qtype:  dns.TypeAAAA,
			Qclass: dns.ClassINET,
		})
	}

	buffer := buf.New()
	common.Must(buffer.Reset(func(b []byte) (int, error) {
		writtenBuffer, err := msg.PackBuffer(b)
		return len(writtenBuffer), err
	}))

	return buffer
}

func (s *UDPNameServer) QueryA(domain string) <-chan *ARecord {
	response := make(chan *ARecord, 1)
	id := s.AssignUnusedID(response)

	ctx, cancel := context.WithCancel(context.Background())
	s.udpServer.Dispatch(ctx, s.address, s.BuildQueryA(domain, id), s.HandleResponse)

	go func() {
		for i := 0; i < 2; i++ {
			time.Sleep(time.Second)
			s.Lock()
			_, found := s.requests[id]
			s.Unlock()
			if !found {
				break
			}
			s.udpServer.Dispatch(ctx, s.address, s.BuildQueryA(domain, id), s.HandleResponse)
		}
		cancel()
	}()

	return response
}

type LocalNameServer struct {
}

func (*LocalNameServer) QueryA(domain string) <-chan *ARecord {
	response := make(chan *ARecord, 1)

	go func() {
		defer close(response)

		ips, err := net.LookupIP(domain)
		if err != nil {
			log.Trace(newError("failed to lookup IPs for domain ", domain).Base(err))
			return
		}

		response <- &ARecord{
			IPs:    ips,
			Expire: time.Now().Add(time.Hour),
		}
	}()

	return response
}
