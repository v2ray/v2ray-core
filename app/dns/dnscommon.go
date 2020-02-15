// +build !confonly

package dns

import (
	"encoding/binary"
	"time"

	"golang.org/x/net/dns/dnsmessage"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	dns_feature "v2ray.com/core/features/dns"
)

// Fqdn normalize domain make sure it ends with '.'
func Fqdn(domain string) string {
	if len(domain) > 0 && domain[len(domain)-1] == '.' {
		return domain
	}
	return domain + "."
}

type record struct {
	A    *IPRecord
	AAAA *IPRecord
}

// IPRecord is a cacheable item for a resolved domain
type IPRecord struct {
	ReqID  uint16
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

func isNewer(baseRec *IPRecord, newRec *IPRecord) bool {
	if newRec == nil {
		return false
	}
	if baseRec == nil {
		return true
	}
	return baseRec.Expire.Before(newRec.Expire)
}

var (
	errRecordNotFound = errors.New("record not found")
)

type dnsRequest struct {
	reqType dnsmessage.Type
	domain  string
	start   time.Time
	expire  time.Time
	msg     *dnsmessage.Message
}

func genEDNS0Options(clientIP net.IP) *dnsmessage.Resource {
	if len(clientIP) == 0 {
		return nil
	}

	var netmask int
	var family uint16

	if len(clientIP) == 4 {
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
		ip := clientIP.To4().Mask(net.CIDRMask(netmask, net.IPv4len*8))
		needLength := (netmask + 8 - 1) / 8 // division rounding up
		b = append(b, ip[:needLength]...)
	case 2:
		ip := clientIP.Mask(net.CIDRMask(netmask, net.IPv6len*8))
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

func buildReqMsgs(domain string, option IPOption, reqIDGen func() uint16, reqOpts *dnsmessage.Resource) []*dnsRequest {
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

	var reqs []*dnsRequest
	now := time.Now()

	if option.IPv4Enable {
		msg := new(dnsmessage.Message)
		msg.Header.ID = reqIDGen()
		msg.Header.RecursionDesired = true
		msg.Questions = []dnsmessage.Question{qA}
		if reqOpts != nil {
			msg.Additionals = append(msg.Additionals, *reqOpts)
		}
		reqs = append(reqs, &dnsRequest{
			reqType: dnsmessage.TypeA,
			domain:  domain,
			start:   now,
			msg:     msg,
		})
	}

	if option.IPv6Enable {
		msg := new(dnsmessage.Message)
		msg.Header.ID = reqIDGen()
		msg.Header.RecursionDesired = true
		msg.Questions = []dnsmessage.Question{qAAAA}
		if reqOpts != nil {
			msg.Additionals = append(msg.Additionals, *reqOpts)
		}
		reqs = append(reqs, &dnsRequest{
			reqType: dnsmessage.TypeAAAA,
			domain:  domain,
			start:   now,
			msg:     msg,
		})
	}

	return reqs
}

// parseResponse parse DNS answers from the returned payload
func parseResponse(payload []byte) (*IPRecord, error) {
	var parser dnsmessage.Parser
	h, err := parser.Start(payload)
	if err != nil {
		return nil, newError("failed to parse DNS response").Base(err).AtWarning()
	}
	if err := parser.SkipAllQuestions(); err != nil {
		return nil, newError("failed to skip questions in DNS response").Base(err).AtWarning()
	}

	now := time.Now()
	ipRecord := &IPRecord{
		ReqID:  h.ID,
		RCode:  h.RCode,
		Expire: now.Add(time.Second * 600),
	}

L:
	for {
		ah, err := parser.AnswerHeader()
		if err != nil {
			if err != dnsmessage.ErrSectionDone {
				newError("failed to parse answer section for domain: ", ah.Name.String()).Base(err).WriteToLog()
			}
			break
		}

		ttl := ah.TTL
		if ttl == 0 {
			ttl = 600
		}
		expire := now.Add(time.Duration(ttl) * time.Second)
		if ipRecord.Expire.After(expire) {
			ipRecord.Expire = expire
		}

		switch ah.Type {
		case dnsmessage.TypeA:
			ans, err := parser.AResource()
			if err != nil {
				newError("failed to parse A record for domain: ", ah.Name).Base(err).WriteToLog()
				break L
			}
			ipRecord.IP = append(ipRecord.IP, net.IPAddress(ans.A[:]))
		case dnsmessage.TypeAAAA:
			ans, err := parser.AAAAResource()
			if err != nil {
				newError("failed to parse A record for domain: ", ah.Name).Base(err).WriteToLog()
				break L
			}
			ipRecord.IP = append(ipRecord.IP, net.IPAddress(ans.AAAA[:]))
		default:
			if err := parser.SkipAnswer(); err != nil {
				newError("failed to skip answer").Base(err).WriteToLog()
				break L
			}
			continue
		}
	}

	return ipRecord, nil
}
