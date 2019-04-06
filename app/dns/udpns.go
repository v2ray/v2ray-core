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
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/dns"
	udp_proto "v2ray.com/core/common/protocol/udp"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal/pubsub"
	"v2ray.com/core/common/task"
	dns_feature "v2ray.com/core/features/dns"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport/internet/udp"
)

type record struct {
	A    *IPRecord
	AAAA *IPRecord
}

type IPRecord struct {
	IP     []net.Address
	Expire time.Time
	RCode  dnsmessage.RCode
}

func (r *IPRecord) getIPs() ([]net.Address, error) {
	if r == nil || r.Expire.Before(time.Now()) {
		return nil, errRecordNotFound
	}
	if r.RCode != dnsmessage.RCodeSuccess {
		return nil, dns_feature.RCodeError(r.RCode)
	}
	return r.IP, nil
}

type pendingRequest struct {
	domain  string
	expire  time.Time
	recType dnsmessage.Type
}

var (
	errRecordNotFound = errors.New("record not found")
)

type ClassicNameServer struct {
	sync.RWMutex
	address   net.Destination
	ips       map[string]record
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
		ips:      make(map[string]record),
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

	for domain, record := range s.ips {
		if record.A != nil && record.A.Expire.Before(now) {
			record.A = nil
		}
		if record.AAAA != nil && record.AAAA.Expire.Before(now) {
			record.AAAA = nil
		}

		if record.A == nil && record.AAAA == nil {
			delete(s.ips, domain)
		} else {
			s.ips[domain] = record
		}
	}

	if len(s.ips) == 0 {
		s.ips = make(map[string]record)
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
	recType := req.recType

	now := time.Now()
	ipRecord := &IPRecord{
		RCode:  header.RCode,
		Expire: now.Add(time.Second * 600),
	}

L:
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
		expire := now.Add(time.Duration(ttl) * time.Second)
		if ipRecord.Expire.After(expire) {
			ipRecord.Expire = expire
		}

		if header.Type != recType {
			if err := parser.SkipAnswer(); err != nil {
				newError("failed to skip answer").Base(err).WriteToLog()
				break L
			}
			continue
		}

		switch header.Type {
		case dnsmessage.TypeA:
			ans, err := parser.AResource()
			if err != nil {
				newError("failed to parse A record for domain: ", domain).Base(err).WriteToLog()
				break L
			}
			ipRecord.IP = append(ipRecord.IP, net.IPAddress(ans.A[:]))
		case dnsmessage.TypeAAAA:
			ans, err := parser.AAAAResource()
			if err != nil {
				newError("failed to parse A record for domain: ", domain).Base(err).WriteToLog()
				break L
			}
			ipRecord.IP = append(ipRecord.IP, net.IPAddress(ans.AAAA[:]))
		default:
			if err := parser.SkipAnswer(); err != nil {
				newError("failed to skip answer").Base(err).WriteToLog()
				break L
			}
		}
	}

	var rec record
	switch recType {
	case dnsmessage.TypeA:
		rec.A = ipRecord
	case dnsmessage.TypeAAAA:
		rec.AAAA = ipRecord
	}

	if len(domain) > 0 && (rec.A != nil || rec.AAAA != nil) {
		s.updateIP(domain, rec)
	}
}

func isNewer(baseRec *IPRecord, newRec *IPRecord) bool {
	if newRec == nil {
		return false
	}
	if baseRec == nil {
		return true
	}
	return baseRec.Expire.Before(newRec.Expire)
}

func (s *ClassicNameServer) updateIP(domain string, newRec record) {
	s.Lock()

	newError("updating IP records for domain:", domain).AtDebug().WriteToLog()
	rec := s.ips[domain]

	updated := false
	if isNewer(rec.A, newRec.A) {
		rec.A = newRec.A
		updated = true
	}
	if isNewer(rec.AAAA, newRec.AAAA) {
		rec.AAAA = newRec.AAAA
		updated = true
	}

	if updated {
		s.ips[domain] = rec
		s.pub.Publish(domain, nil)
	}

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

func (s *ClassicNameServer) addPendingRequest(domain string, recType dnsmessage.Type) uint16 {
	id := uint16(atomic.AddUint32(&s.reqID, 1))
	s.Lock()
	defer s.Unlock()

	s.requests[id] = pendingRequest{
		domain:  domain,
		expire:  time.Now().Add(time.Second * 8),
		recType: recType,
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
		msg.Header.ID = s.addPendingRequest(domain, dnsmessage.TypeA)
		msg.Header.RecursionDesired = true
		msg.Questions = []dnsmessage.Question{qA}
		if opt := s.getMsgOptions(); opt != nil {
			msg.Additionals = append(msg.Additionals, *opt)
		}
		msgs = append(msgs, msg)
	}

	if option.IPv6Enable {
		msg := new(dnsmessage.Message)
		msg.Header.ID = s.addPendingRequest(domain, dnsmessage.TypeAAAA)
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
		udpCtx = session.ContextWithContent(udpCtx, &session.Content{
			Protocol: "dns",
		})
		s.udpServer.Dispatch(udpCtx, s.address, b)
	}
}

func (s *ClassicNameServer) findIPsForDomain(domain string, option IPOption) ([]net.IP, error) {
	s.RLock()
	record, found := s.ips[domain]
	s.RUnlock()

	if !found {
		return nil, errRecordNotFound
	}

	var ips []net.Address
	var lastErr error
	if option.IPv4Enable {
		a, err := record.A.getIPs()
		if err != nil {
			lastErr = err
		}
		ips = append(ips, a...)
	}

	if option.IPv6Enable {
		aaaa, err := record.AAAA.getIPs()
		if err != nil {
			lastErr = err
		}
		ips = append(ips, aaaa...)
	}

	if len(ips) > 0 {
		return toNetIP(ips), nil
	}

	if lastErr != nil {
		return nil, lastErr
	}

	return nil, dns_feature.ErrEmptyResponse
}

func Fqdn(domain string) string {
	if len(domain) > 0 && domain[len(domain)-1] == '.' {
		return domain
	}
	return domain + "."
}

func (s *ClassicNameServer) QueryIP(ctx context.Context, domain string, option IPOption) ([]net.IP, error) {
	fqdn := Fqdn(domain)

	ips, err := s.findIPsForDomain(fqdn, option)
	if err != errRecordNotFound {
		return ips, err
	}

	sub := s.pub.Subscribe(fqdn)
	defer sub.Close()

	s.sendQuery(ctx, fqdn, option)

	for {
		ips, err := s.findIPsForDomain(fqdn, option)
		if err != errRecordNotFound {
			return ips, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-sub.Wait():
		}
	}
}
