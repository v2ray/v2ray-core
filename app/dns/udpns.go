// +build !confonly

package dns

import (
	"context"
	"strings"
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
	dns_feature "v2ray.com/core/features/dns"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport/internet/udp"
)

type ClassicNameServer struct {
	sync.RWMutex
	name      string
	address   net.Destination
	ips       map[string]record
	requests  map[uint16]dnsRequest
	pub       *pubsub.Service
	udpServer *udp.Dispatcher
	cleanup   *task.Periodic
	reqID     uint32
	clientIP  net.IP
}

func NewClassicNameServer(address net.Destination, dispatcher routing.Dispatcher, clientIP net.IP) *ClassicNameServer {

	// default to 53 if unspecific
	if address.Port == 0 {
		address.Port = net.Port(53)
	}

	s := &ClassicNameServer{
		address:  address,
		ips:      make(map[string]record),
		requests: make(map[uint16]dnsRequest),
		clientIP: clientIP,
		pub:      pubsub.NewService(),
		name:     strings.ToUpper(address.String()),
	}
	s.cleanup = &task.Periodic{
		Interval: time.Minute,
		Execute:  s.Cleanup,
	}
	s.udpServer = udp.NewDispatcher(dispatcher, s.HandleResponse)
	newError("DNS: created udp client inited for ", address.NetAddr()).AtInfo().WriteToLog()
	return s
}

func (s *ClassicNameServer) Name() string {
	return s.name
}

func (s *ClassicNameServer) Cleanup() error {
	now := time.Now()
	s.Lock()
	defer s.Unlock()

	if len(s.ips) == 0 && len(s.requests) == 0 {
		return newError(s.name, " nothing to do. stopping...")
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
		s.requests = make(map[uint16]dnsRequest)
	}

	return nil
}

func (s *ClassicNameServer) HandleResponse(ctx context.Context, packet *udp_proto.Packet) {

	ipRec, err := parseResponse(packet.Payload.Bytes())
	if err != nil {
		newError(s.name, " fail to parse responsed DNS udp").AtError().WriteToLog()
		return
	}

	s.Lock()
	id := ipRec.ReqID
	req, ok := s.requests[id]
	if ok {
		// remove the pending request
		delete(s.requests, id)
	}
	s.Unlock()
	if !ok {
		newError(s.name, " cannot find the pending request").AtError().WriteToLog()
		return
	}

	var rec record
	switch req.reqType {
	case dnsmessage.TypeA:
		rec.A = ipRec
	case dnsmessage.TypeAAAA:
		rec.AAAA = ipRec
	}

	elapsed := time.Since(req.start)
	newError(s.name, " got answere: ", req.domain, " ", req.reqType, " -> ", ipRec.IP, " ", elapsed).AtInfo().WriteToLog()
	if len(req.domain) > 0 && (rec.A != nil || rec.AAAA != nil) {
		s.updateIP(req.domain, rec)
	}
}

func (s *ClassicNameServer) updateIP(domain string, newRec record) {
	s.Lock()

	newError(s.name, " updating IP records for domain:", domain).AtDebug().WriteToLog()
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

func (s *ClassicNameServer) newReqID() uint16 {
	return uint16(atomic.AddUint32(&s.reqID, 1))
}

func (s *ClassicNameServer) addPendingRequest(req *dnsRequest) {
	s.Lock()
	defer s.Unlock()

	id := req.msg.ID
	req.expire = time.Now().Add(time.Second * 8)
	s.requests[id] = *req
}

func (s *ClassicNameServer) sendQuery(ctx context.Context, domain string, option IPOption) {
	newError(s.name, " querying DNS for: ", domain).AtDebug().WriteToLog(session.ExportIDToError(ctx))

	reqs := buildReqMsgs(domain, option, s.newReqID, genEDNS0Options(s.clientIP))

	for _, req := range reqs {
		s.addPendingRequest(req)
		b, _ := dns.PackMessage(req.msg)
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

func (s *ClassicNameServer) QueryIP(ctx context.Context, domain string, option IPOption) ([]net.IP, error) {

	fqdn := Fqdn(domain)

	ips, err := s.findIPsForDomain(fqdn, option)
	if err != errRecordNotFound {
		newError(s.name, " cache HIT ", domain, " -> ", ips).Base(err).AtDebug().WriteToLog()
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
