// +build !confonly

package dns

import (
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/dns/dnsmessage"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/dns"
	udp_proto "v2ray.com/core/common/protocol/udp"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal/pubsub"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport/internet/udp"
)

type IPRecord struct {
	IP     net.Address
	Expire time.Time
}

type pendingRequest struct {
	domain string
	expire time.Time
}

type ClassicNameServer struct {
	sync.RWMutex
	address   net.Destination
	ips       map[string][]IPRecord
	requests  map[uint16]pendingRequest
	pub       *pubsub.Service
	udpServer *udp.Dispatcher
	cleanup   *task.Periodic
	reqID     uint32
	clientIP  net.IP
}

func NewClassicNameServer(address net.Destination, dispatcher routing.Dispatcher, clientIP net.IP) *ClassicNameServer {
	s := &ClassicNameServer{
		address:  address,
		ips:      make(map[string][]IPRecord),
		requests: make(map[uint16]pendingRequest),
		clientIP: clientIP,
		pub:      pubsub.NewService(),
	}
	s.cleanup = &task.Periodic{
		Interval: time.Minute,
		Execute:  s.Cleanup,
	}
	s.udpServer = udp.NewDispatcher(dispatcher, s.HandleResponse)
	return s
}

func (s *ClassicNameServer) Name() string {
	return s.address.String()
}

func (s *ClassicNameServer) Cleanup() error {
	now := time.Now()
	s.Lock()
	defer s.Unlock()

	if len(s.ips) == 0 && len(s.requests) == 0 {
		return newError("nothing to do. stopping...")
	}

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

	for id, req := range s.requests {
		if req.expire.Before(now) {
			delete(s.requests, id)
		}
	}

	if len(s.requests) == 0 {
		s.requests = make(map[uint16]pendingRequest)
	}

	return nil
}

func (s *ClassicNameServer) HandleResponse(ctx context.Context, packet *udp_proto.Packet) {
	payload := packet.Payload

	var parser dnsmessage.Parser
	header, err := parser.Start(payload.Bytes())
	if err != nil {
		newError("failed to parse DNS response").Base(err).AtWarning().WriteToLog()
		return
	}
	if err := parser.SkipAllQuestions(); err != nil {
		newError("failed to skip questions in DNS response").Base(err).AtWarning().WriteToLog()
		return
	}

	id := header.ID
	s.Lock()
	req, f := s.requests[id]
	if f {
		delete(s.requests, id)
	}
	s.Unlock()

	if !f {
		return
	}

	domain := req.domain
	ips := make([]IPRecord, 0, 16)

	now := time.Now()
	for {
		header, err := parser.AnswerHeader()
		if err != nil {
			if err != dnsmessage.ErrSectionDone {
				newError("failed to parse answer section for domain: ", domain).Base(err).WriteToLog()
			}
			break
		}
		ttl := header.TTL
		if ttl == 0 {
			ttl = 600
		}
		switch header.Type {
		case dnsmessage.TypeA:
			ans, err := parser.AResource()
			if err != nil {
				newError("failed to parse A record for domain: ", domain).Base(err).WriteToLog()
				break
			}
			ips = append(ips, IPRecord{
				IP:     net.IPAddress(ans.A[:]),
				Expire: now.Add(time.Duration(ttl) * time.Second),
			})
		case dnsmessage.TypeAAAA:
			ans, err := parser.AAAAResource()
			if err != nil {
				newError("failed to parse A record for domain: ", domain).Base(err).WriteToLog()
				break
			}
			ips = append(ips, IPRecord{
				IP:     net.IPAddress(ans.AAAA[:]),
				Expire: now.Add(time.Duration(ttl) * time.Second),
			})
		default:
			if err := parser.SkipAnswer(); err != nil {
				newError("failed to skip answer").Base(err).WriteToLog()
			}
		}
	}

	if len(domain) > 0 && len(ips) > 0 {
		s.updateIP(domain, ips)
	}
}

func (s *ClassicNameServer) updateIP(domain string, ips []IPRecord) {
	s.Lock()

	newError("updating IP records for domain:", domain).AtDebug().WriteToLog()
	now := time.Now()
	eips := s.ips[domain]
	for _, ip := range eips {
		if ip.Expire.After(now) {
			ips = append(ips, ip)
		}
	}
	s.ips[domain] = ips
	s.pub.Publish(domain, nil)

	s.Unlock()
	common.Must(s.cleanup.Start())
}

func (s *ClassicNameServer) getMsgOptions() *dnsmessage.Resource {
	if len(s.clientIP) == 0 {
		return nil
	}

	var netmask int
	var family uint16

	if len(s.clientIP) == 4 {
		family = 1
		netmask = 24 // 24 for IPV4, 96 for IPv6
	} else {
		family = 2
		netmask = 96
	}

	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[0:], family)
	b[2] = byte(netmask)
	b[3] = 0
	switch family {
	case 1:
		ip := s.clientIP.To4().Mask(net.CIDRMask(netmask, net.IPv4len*8))
		needLength := (netmask + 8 - 1) / 8 // division rounding up
		b = append(b, ip[:needLength]...)
	case 2:
		ip := s.clientIP.Mask(net.CIDRMask(netmask, net.IPv6len*8))
		needLength := (netmask + 8 - 1) / 8 // division rounding up
		b = append(b, ip[:needLength]...)
	}

	const EDNS0SUBNET = 0x08

	opt := new(dnsmessage.Resource)
	common.Must(opt.Header.SetEDNS0(1350, 0xfe00, true))

	opt.Body = &dnsmessage.OPTResource{
		Options: []dnsmessage.Option{
			{
				Code: EDNS0SUBNET,
				Data: b,
			},
		},
	}

	return opt
}

func (s *ClassicNameServer) addPendingRequest(domain string) uint16 {
	id := uint16(atomic.AddUint32(&s.reqID, 1))
	s.Lock()
	defer s.Unlock()

	s.requests[id] = pendingRequest{
		domain: domain,
		expire: time.Now().Add(time.Second * 8),
	}

	return id
}

func (s *ClassicNameServer) buildMsgs(domain string, option IPOption) []*dnsmessage.Message {
	qA := dnsmessage.Question{
		Name:  dnsmessage.MustNewName(domain),
		Type:  dnsmessage.TypeA,
		Class: dnsmessage.ClassINET,
	}

	qAAAA := dnsmessage.Question{
		Name:  dnsmessage.MustNewName(domain),
		Type:  dnsmessage.TypeAAAA,
		Class: dnsmessage.ClassINET,
	}

	var msgs []*dnsmessage.Message

	if option.IPv4Enable {
		msg := new(dnsmessage.Message)
		msg.Header.ID = s.addPendingRequest(domain)
		msg.Header.RecursionDesired = true
		msg.Questions = []dnsmessage.Question{qA}
		if opt := s.getMsgOptions(); opt != nil {
			msg.Additionals = append(msg.Additionals, *opt)
		}
		msgs = append(msgs, msg)
	}

	if option.IPv6Enable {
		msg := new(dnsmessage.Message)
		msg.Header.ID = s.addPendingRequest(domain)
		msg.Header.RecursionDesired = true
		msg.Questions = []dnsmessage.Question{qAAAA}
		if opt := s.getMsgOptions(); opt != nil {
			msg.Additionals = append(msg.Additionals, *opt)
		}
		msgs = append(msgs, msg)
	}

	return msgs
}

func (s *ClassicNameServer) sendQuery(ctx context.Context, domain string, option IPOption) {
	newError("querying DNS for: ", domain).AtDebug().WriteToLog(session.ExportIDToError(ctx))

	msgs := s.buildMsgs(domain, option)

	for _, msg := range msgs {
		b, err := dns.PackMessage(msg)
		common.Must(err)
		udpCtx := context.Background()
		if inbound := session.InboundFromContext(ctx); inbound != nil {
			udpCtx = session.ContextWithInbound(udpCtx, inbound)
		}
		s.udpServer.Dispatch(udpCtx, s.address, b)
	}
}

func (s *ClassicNameServer) findIPsForDomain(domain string, option IPOption) []net.IP {
	s.RLock()
	records, found := s.ips[domain]
	s.RUnlock()

	if found && len(records) > 0 {
		var ips []net.Address
		now := time.Now()
		for _, rec := range records {
			if rec.Expire.After(now) {
				ips = append(ips, rec.IP)
			}
		}
		return toNetIP(filterIP(ips, option))
	}
	return nil
}

func Fqdn(domain string) string {
	if len(domain) > 0 && domain[len(domain)-1] == '.' {
		return domain
	}
	return domain + "."
}

func (s *ClassicNameServer) QueryIP(ctx context.Context, domain string, option IPOption) ([]net.IP, error) {
	fqdn := Fqdn(domain)

	ips := s.findIPsForDomain(fqdn, option)
	if len(ips) > 0 {
		return ips, nil
	}

	sub := s.pub.Subscribe(fqdn)
	defer sub.Close()

	s.sendQuery(ctx, fqdn, option)

	for {
		ips := s.findIPsForDomain(fqdn, option)
		if len(ips) > 0 {
			return ips, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-sub.Wait():
		}
	}
}
