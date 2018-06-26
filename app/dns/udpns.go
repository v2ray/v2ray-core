package dns

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/transport/internet/udp"
)

var (
	multiQuestionDNS = map[net.Address]bool{
		net.IPAddress([]byte{8, 8, 8, 8}): true,
		net.IPAddress([]byte{8, 8, 4, 4}): true,
		net.IPAddress([]byte{9, 9, 9, 9}): true,
	}
)

type IPRecord struct {
	IP     net.IP
	Expire time.Time
}

type ClassicNameServer struct {
	sync.RWMutex
	address   net.Destination
	ips       map[string][]IPRecord
	updated   signal.Notifier
	udpServer *udp.Dispatcher
	cleanup   *task.Periodic
	reqID     uint32
	clientIP  net.IP
}

func NewClassicNameServer(address net.Destination, dispatcher core.Dispatcher, clientIP net.IP) *ClassicNameServer {
	s := &ClassicNameServer{
		address:   address,
		ips:       make(map[string][]IPRecord),
		udpServer: udp.NewDispatcher(dispatcher),
		clientIP:  clientIP,
	}
	s.cleanup = &task.Periodic{
		Interval: time.Minute,
		Execute:  s.Cleanup,
	}
	common.Must(s.cleanup.Start())
	return s
}

func (s *ClassicNameServer) Cleanup() error {
	now := time.Now()
	s.Lock()
	for domain, ips := range s.ips {
		newIPs := make([]IPRecord, 0, len(ips))
		for _, ip := range ips {
			if ip.Expire.After(now) {
				newIPs = append(newIPs, ip)
			}
		}
		if len(newIPs) == 0 {
			delete(s.ips, domain)
		} else if len(newIPs) < len(ips) {
			s.ips[domain] = newIPs
		}
	}

	if len(s.ips) == 0 {
		s.ips = make(map[string][]IPRecord)
	}

	s.Unlock()
	return nil
}

func (s *ClassicNameServer) HandleResponse(payload *buf.Buffer) {
	msg := new(dns.Msg)
	err := msg.Unpack(payload.Bytes())
	if err == dns.ErrTruncated {
		newError("truncated message received. DNS server should still work. If you see anything abnormal, please submit an issue to v2ray-core.").AtWarning().WriteToLog()
	} else if err != nil {
		newError("failed to parse DNS response").Base(err).AtWarning().WriteToLog()
		return
	}

	var domain string
	ips := make([]IPRecord, 0, 16)

	now := time.Now()
	for _, rr := range msg.Answer {
		var ip net.IP
		domain = rr.Header().Name
		ttl := rr.Header().Ttl
		switch rr := rr.(type) {
		case *dns.A:
			ip = rr.A
		case *dns.AAAA:
			ip = rr.AAAA
		}
		if ttl == 0 {
			ttl = 300
		}
		if len(ip) > 0 {
			ips = append(ips, IPRecord{
				IP:     ip,
				Expire: now.Add(time.Second * time.Duration(ttl)),
			})
		}
	}

	if len(domain) > 0 && len(ips) > 0 {
		s.updateIP(domain, ips)
	}
}

func (s *ClassicNameServer) updateIP(domain string, ips []IPRecord) {
	s.Lock()
	defer s.Unlock()

	newError("updating IP records for domain:", domain).AtDebug().WriteToLog()
	now := time.Now()
	eips := s.ips[domain]
	for _, ip := range eips {
		if ip.Expire.After(now) {
			ips = append(ips, ip)
		}
	}
	s.ips[domain] = ips
	s.updated.Signal()
}

func (s *ClassicNameServer) getMsgOptions() *dns.OPT {
	if len(s.clientIP) == 0 {
		return nil
	}

	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	o.SetUDPSize(1280)

	e := new(dns.EDNS0_SUBNET)
	e.Code = dns.EDNS0SUBNET
	if len(s.clientIP) == 4 {
		e.Family = 1 // 1 for IPv4 source address, 2 for IPv6
	} else {
		e.Family = 2
	}

	e.SourceNetmask = 24 // 32 for IPV4, 128 for IPv6
	e.SourceScope = 0
	e.Address = s.clientIP
	o.Option = append(o.Option, e)

	return o

}

func (s *ClassicNameServer) buildMsgs(domain string) []*dns.Msg {
	allowMulti := multiQuestionDNS[s.address.Address]

	qA := dns.Question{
		Name:   domain,
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}

	qAAAA := dns.Question{
		Name:   domain,
		Qtype:  dns.TypeAAAA,
		Qclass: dns.ClassINET,
	}

	var msgs []*dns.Msg

	{
		msg := new(dns.Msg)
		msg.Id = uint16(atomic.AddUint32(&s.reqID, 1))
		msg.RecursionDesired = true
		msg.Question = []dns.Question{qA}
		if allowMulti {
			msg.Question = append(msg.Question, qAAAA)
		}
		if opt := s.getMsgOptions(); opt != nil {
			msg.Extra = append(msg.Extra, opt)
		}
		msgs = append(msgs, msg)
	}

	if !allowMulti {
		msg := new(dns.Msg)
		msg.Id = uint16(atomic.AddUint32(&s.reqID, 1))
		msg.RecursionDesired = true
		msg.Question = []dns.Question{qAAAA}
		if opt := s.getMsgOptions(); opt != nil {
			msg.Extra = append(msg.Extra, opt)
		}
		msgs = append(msgs, msg)
	}

	return msgs
}

func msgToBuffer(msg *dns.Msg) (*buf.Buffer, error) {
	buffer := buf.New()
	if err := buffer.Reset(func(b []byte) (int, error) {
		writtenBuffer, err := msg.PackBuffer(b)
		return len(writtenBuffer), err
	}); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (s *ClassicNameServer) sendQuery(ctx context.Context, domain string) {
	msgs := s.buildMsgs(domain)

	for _, msg := range msgs {
		b, err := msgToBuffer(msg)
		common.Must(err)
		s.udpServer.Dispatch(ctx, s.address, b, s.HandleResponse)
	}
}

func (s *ClassicNameServer) findIPsForDomain(domain string) []net.IP {
	records, found := s.ips[domain]
	if found && len(records) > 0 {
		var ips []net.IP
		now := time.Now()
		for _, rec := range records {
			if rec.Expire.After(now) {
				ips = append(ips, rec.IP)
			}
		}
		return ips
	}
	return nil
}

func (s *ClassicNameServer) QueryIP(ctx context.Context, domain string) ([]net.IP, error) {
	fqdn := dns.Fqdn(domain)

	ips := s.findIPsForDomain(fqdn)
	if len(ips) > 0 {
		return ips, nil
	}

	s.sendQuery(ctx, fqdn)

	for {
		ips := s.findIPsForDomain(fqdn)
		if len(ips) > 0 {
			return ips, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-s.updated.Wait():
		}
	}
}
