// +build !confonly

package dns

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
	dns_feature "v2ray.com/core/features/dns"

	"golang.org/x/net/dns/dnsmessage"
	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol/dns"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal/pubsub"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/routing"
)

// DoHNameServer implimented DNS over HTTPS (RFC8484) Wire Format,
// which is compatiable with traditional dns over udp(RFC1035),
// thus most of the DOH implimentation is copied from udpns.go
type DoHNameServer struct {
	sync.RWMutex
	dispatcher routing.Dispatcher
	dohDests   []net.Destination
	ips        map[string]record
	pub        *pubsub.Service
	cleanup    *task.Periodic
	reqID      uint32
	clientIP   net.IP
	httpClient *http.Client
	dohURL     string
	name       string
}

// NewDoHNameServer creates DOH client object for remote resolving
func NewDoHNameServer(url *url.URL, dispatcher routing.Dispatcher, clientIP net.IP) (*DoHNameServer, error) {

	dohAddr := net.ParseAddress(url.Hostname())
	dohPort := "443"
	if url.Port() != "" {
		dohPort = url.Port()
	}

	parseIPDest := func(ip net.IP, port string) net.Destination {
		strIP := ip.String()
		if len(ip) == net.IPv6len {
			strIP = fmt.Sprintf("[%s]", strIP)
		}
		dest, err := net.ParseDestination(fmt.Sprintf("tcp:%s:%s", strIP, port))
		common.Must(err)
		return dest
	}

	var dests []net.Destination
	if dohAddr.Family().IsDomain() {
		// resolve DOH server in advance
		ips, err := net.LookupIP(dohAddr.Domain())
		if err != nil || len(ips) == 0 {
			return nil, err
		}
		for _, ip := range ips {
			dests = append(dests, parseIPDest(ip, dohPort))
		}
	} else {
		ip := dohAddr.IP()
		dests = append(dests, parseIPDest(ip, dohPort))
	}

	newError("DNS: created Remote DOH client for ", url.String(), ", preresolved Dests: ", dests).AtInfo().WriteToLog()
	s := baseDOHNameServer(url, "DOH", clientIP)
	s.dispatcher = dispatcher
	s.dohDests = dests

	// Dispatched connection will be closed (interupted) after each request
	// This makes DOH inefficient without a keeped-alive connection
	// See: core/app/proxyman/outbound/handler.go:113
	// Using mux (https request wrapped in a stream layer) improves the situation.
	// Recommand to use NewDoHLocalNameServer (DOHL:) if v2ray instance is running on
	//  a normal network eg. the server side of v2ray
	tr := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DialContext:         s.DialContext,
	}

	dispatchedClient := &http.Client{
		Transport: tr,
		Timeout:   16 * time.Second,
	}

	s.httpClient = dispatchedClient
	return s, nil
}

// NewDoHLocalNameServer creates DOH client object for local resolving
func NewDoHLocalNameServer(url *url.URL, clientIP net.IP) *DoHNameServer {
	url.Scheme = "https"
	s := baseDOHNameServer(url, "DOHL", clientIP)
	s.httpClient = &http.Client{
		Timeout: time.Second * 180,
	}
	newError("DNS: created Local DOH client for ", url.String()).AtInfo().WriteToLog()
	return s
}

func baseDOHNameServer(url *url.URL, prefix string, clientIP net.IP) *DoHNameServer {

	s := &DoHNameServer{
		ips:      make(map[string]record),
		clientIP: clientIP,
		pub:      pubsub.NewService(),
		name:     fmt.Sprintf("%s//%s", prefix, url.Host),
		dohURL:   url.String(),
	}
	s.cleanup = &task.Periodic{
		Interval: time.Minute,
		Execute:  s.Cleanup,
	}

	return s
}

// Name returns client name
func (s *DoHNameServer) Name() string {
	return s.name
}

// DialContext offer dispatched connection through core routing
func (s *DoHNameServer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {

	dest := s.dohDests[dice.Roll(len(s.dohDests))]

	link, err := s.dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return nil, err
	}
	return net.NewConnection(
		net.ConnectionInputMulti(link.Writer),
		net.ConnectionOutputMulti(link.Reader),
	), nil
}

// Cleanup clears expired items from cache
func (s *DoHNameServer) Cleanup() error {
	now := time.Now()
	s.Lock()
	defer s.Unlock()

	if len(s.ips) == 0 {
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
			newError(s.name, " cleanup ", domain).AtDebug().WriteToLog()
			delete(s.ips, domain)
		} else {
			s.ips[domain] = record
		}
	}

	if len(s.ips) == 0 {
		s.ips = make(map[string]record)
	}

	return nil
}

func (s *DoHNameServer) updateIP(req *dnsRequest, ipRec *IPRecord) {
	elapsed := time.Since(req.start)

	s.Lock()
	rec := s.ips[req.domain]
	updated := false

	switch req.reqType {
	case dnsmessage.TypeA:
		if isNewer(rec.A, ipRec) {
			rec.A = ipRec
			updated = true
		}
	case dnsmessage.TypeAAAA:
		addr := make([]net.Address, 0)
		for _, ip := range ipRec.IP {
			if len(ip.IP()) == net.IPv6len {
				addr = append(addr, ip)
			}
		}
		ipRec.IP = addr
		if isNewer(rec.AAAA, ipRec) {
			rec.AAAA = ipRec
			updated = true
		}
	}
	newError(s.name, " got answere: ", req.domain, " ", req.reqType, " -> ", ipRec.IP, " ", elapsed).AtInfo().WriteToLog()

	if updated {
		s.ips[req.domain] = rec
	}
	switch req.reqType {
	case dnsmessage.TypeA:
		s.pub.Publish(req.domain+"4", nil)
	case dnsmessage.TypeAAAA:
		s.pub.Publish(req.domain+"6", nil)
	}
	s.Unlock()
	common.Must(s.cleanup.Start())
}

func (s *DoHNameServer) newReqID() uint16 {
	return uint16(atomic.AddUint32(&s.reqID, 1))
}

func (s *DoHNameServer) sendQuery(ctx context.Context, domain string, option IPOption) {
	newError(s.name, " querying: ", domain).AtInfo().WriteToLog(session.ExportIDToError(ctx))

	reqs := buildReqMsgs(domain, option, s.newReqID, genEDNS0Options(s.clientIP))

	var deadline time.Time
	if d, ok := ctx.Deadline(); ok {
		deadline = d
	} else {
		deadline = time.Now().Add(time.Second * 8)
	}

	for _, req := range reqs {

		go func(r *dnsRequest) {

			// generate new context for each req, using same context
			// may cause reqs all aborted if any one encounter an error
			dnsCtx := context.Background()

			// reserve internal dns server requested Inbound
			if inbound := session.InboundFromContext(ctx); inbound != nil {
				dnsCtx = session.ContextWithInbound(dnsCtx, inbound)
			}

			dnsCtx = session.ContextWithContent(dnsCtx, &session.Content{
				Protocol: "https",
			})

			// forced to use mux for DOH
			dnsCtx = session.ContextWithMuxPrefered(dnsCtx, true)

			dnsCtx, cancel := context.WithDeadline(dnsCtx, deadline)
			defer cancel()

			b, _ := dns.PackMessage(r.msg)
			resp, err := s.dohHTTPSContext(dnsCtx, b.Bytes())
			if err != nil {
				newError("failed to retrive response").Base(err).AtError().WriteToLog()
				return
			}
			rec, err := parseResponse(resp)
			if err != nil {
				newError("failed to handle DOH response").Base(err).AtError().WriteToLog()
				return
			}
			s.updateIP(r, rec)
		}(req)
	}
}

func (s *DoHNameServer) dohHTTPSContext(ctx context.Context, b []byte) ([]byte, error) {

	body := bytes.NewBuffer(b)
	req, err := http.NewRequest("POST", s.dohURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/dns-message")
	req.Header.Add("Content-Type", "application/dns-message")

	resp, err := s.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("DOH HTTPS server returned with non-OK code %d", resp.StatusCode)
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *DoHNameServer) findIPsForDomain(domain string, option IPOption) ([]net.IP, error) {
	s.RLock()
	record, found := s.ips[domain]
	s.RUnlock()

	if !found {
		return nil, errRecordNotFound
	}

	var ips []net.Address
	var lastErr error
	if option.IPv6Enable && record.AAAA != nil && record.AAAA.RCode == dnsmessage.RCodeSuccess {
		aaaa, err := record.AAAA.getIPs()
		if err != nil {
			lastErr = err
		}
		ips = append(ips, aaaa...)
	}

	if option.IPv4Enable && record.A != nil && record.A.RCode == dnsmessage.RCodeSuccess {
		a, err := record.A.getIPs()
		if err != nil {
			lastErr = err
		}
		ips = append(ips, a...)
	}

	if len(ips) > 0 {
		return toNetIP(ips), nil
	}

	if lastErr != nil {
		return nil, lastErr
	}

	if (option.IPv4Enable && record.A != nil) || (option.IPv6Enable && record.AAAA != nil) {
		return nil, dns_feature.ErrEmptyResponse
	}

	return nil, errRecordNotFound
}

// QueryIP is called from dns.Server->queryIPTimeout
func (s *DoHNameServer) QueryIP(ctx context.Context, domain string, option IPOption) ([]net.IP, error) {
	fqdn := Fqdn(domain)

	ips, err := s.findIPsForDomain(fqdn, option)
	if err != errRecordNotFound {
		newError(s.name, " cache HIT ", domain, " -> ", ips).Base(err).AtDebug().WriteToLog()
		return ips, err
	}

	// ipv4 and ipv6 belong to different subscription groups
	var sub4, sub6 *pubsub.Subscriber
	if option.IPv4Enable {
		sub4 = s.pub.Subscribe(fqdn + "4")
		defer sub4.Close()
	}
	if option.IPv6Enable {
		sub6 = s.pub.Subscribe(fqdn + "6")
		defer sub6.Close()
	}
	done := make(chan interface{})
	go func() {
		if sub4 != nil {
			select {
			case <-sub4.Wait():
			case <-ctx.Done():
			}
		}
		if sub6 != nil {
			select {
			case <-sub6.Wait():
			case <-ctx.Done():
			}
		}
		close(done)
	}()
	s.sendQuery(ctx, fqdn, option)

	for {
		ips, err := s.findIPsForDomain(fqdn, option)
		if err != errRecordNotFound {
			return ips, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-done:
		}
	}
}
