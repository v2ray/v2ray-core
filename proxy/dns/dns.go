// +build !confonly

package dns

import (
	"context"
	"io"
	"sync"

	"golang.org/x/net/dns/dnsmessage"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	dns_proto "v2ray.com/core/common/protocol/dns"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		h := new(Handler)
		if err := core.RequireFeatures(ctx, func(dnsClient dns.Client) error {
			return h.Init(config.(*Config), dnsClient)
		}); err != nil {
			return nil, err
		}
		return h, nil
	}))
}

type ownLinkVerifier interface {
	IsOwnLink(ctx context.Context) bool
}

type Handler struct {
	ipv4Lookup      dns.IPv4Lookup
	ipv6Lookup      dns.IPv6Lookup
	ownLinkVerifier ownLinkVerifier
}

func (h *Handler) Init(config *Config, dnsClient dns.Client) error {
	ipv4lookup, ok := dnsClient.(dns.IPv4Lookup)
	if !ok {
		return newError("dns.Client doesn't implement IPv4Lookup")
	}
	h.ipv4Lookup = ipv4lookup

	ipv6lookup, ok := dnsClient.(dns.IPv6Lookup)
	if !ok {
		return newError("dns.Client doesn't implement IPv6Lookup")
	}
	h.ipv6Lookup = ipv6lookup

	if v, ok := dnsClient.(ownLinkVerifier); ok {
		h.ownLinkVerifier = v
	}
	return nil
}

func (h *Handler) isOwnLink(ctx context.Context) bool {
	return h.ownLinkVerifier != nil && h.ownLinkVerifier.IsOwnLink(ctx)
}

func parseIPQuery(b []byte) (r bool, domain string, id uint16, qType dnsmessage.Type) {
	var parser dnsmessage.Parser
	header, err := parser.Start(b)
	if err != nil {
		newError("parser start").Base(err).WriteToLog()
		return
	}

	id = header.ID
	q, err := parser.Question()
	if err != nil {
		newError("question").Base(err).WriteToLog()
		return
	}
	qType = q.Type
	if qType != dnsmessage.TypeA && qType != dnsmessage.TypeAAAA {
		return
	}

	domain = q.Name.String()
	r = true
	return
}

// Process implements proxy.Outbound.
func (h *Handler) Process(ctx context.Context, link *transport.Link, d internet.Dialer) error {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return newError("invalid outbound")
	}

	dest := outbound.Target

	conn := &outboundConn{
		dialer: func() (internet.Connection, error) {
			return d.Dial(ctx, dest)
		},
		connReady: make(chan struct{}, 1),
	}

	var reader dns_proto.MessageReader
	var writer dns_proto.MessageWriter
	if dest.Network == net.Network_TCP {
		reader = dns_proto.NewTCPReader(link.Reader)
		writer = &dns_proto.TCPWriter{
			Writer: link.Writer,
		}
	} else {
		reader = &dns_proto.UDPReader{
			Reader: link.Reader,
		}
		writer = &dns_proto.UDPWriter{
			Writer: link.Writer,
		}
	}

	var connReader dns_proto.MessageReader
	var connWriter dns_proto.MessageWriter
	if dest.Network == net.Network_TCP {
		connReader = dns_proto.NewTCPReader(buf.NewReader(conn))
		connWriter = &dns_proto.TCPWriter{
			Writer: buf.NewWriter(conn),
		}
	} else {
		connReader = &dns_proto.UDPReader{
			Reader: &buf.PacketReader{Reader: conn},
		}
		connWriter = &dns_proto.UDPWriter{
			Writer: buf.NewWriter(conn),
		}
	}

	request := func() error {
		defer conn.Close()

		for {
			b, err := reader.ReadMessage()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			if !h.isOwnLink(ctx) {
				isIPQuery, domain, id, qType := parseIPQuery(b.Bytes())
				if isIPQuery {
					go h.handleIPQuery(id, qType, domain, writer)
					continue
				}
			}

			if err := connWriter.WriteMessage(b); err != nil {
				return err
			}
		}
	}

	response := func() error {
		for {
			b, err := connReader.ReadMessage()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			if err := writer.WriteMessage(b); err != nil {
				return err
			}
		}
	}

	if err := task.Run(ctx, request, response); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

func (h *Handler) handleIPQuery(id uint16, qType dnsmessage.Type, domain string, writer dns_proto.MessageWriter) {
	var ips []net.IP
	var err error

	switch qType {
	case dnsmessage.TypeA:
		ips, err = h.ipv4Lookup.LookupIPv4(domain)
	case dnsmessage.TypeAAAA:
		ips, err = h.ipv6Lookup.LookupIPv6(domain)
	}

	if err != nil {
		newError("ip query").Base(err).WriteToLog()
		return
	}

	if len(ips) == 0 {
		return
	}

	b := buf.New()
	rawBytes := b.Extend(buf.Size)
	builder := dnsmessage.NewBuilder(rawBytes[:0], dnsmessage.Header{
		ID:       id,
		RCode:    dnsmessage.RCodeSuccess,
		Response: true,
	})
	builder.StartAnswers()

	rHeader := dnsmessage.ResourceHeader{Name: dnsmessage.MustNewName(domain), Class: dnsmessage.ClassINET, TTL: 600}
	for _, ip := range ips {
		if len(ip) == net.IPv4len {
			var r dnsmessage.AResource
			copy(r.A[:], ip)
			builder.AResource(rHeader, r)
		} else {
			var r dnsmessage.AAAAResource
			copy(r.AAAA[:], ip)
			builder.AAAAResource(rHeader, r)
		}
	}
	msgBytes, err := builder.Finish()
	if err != nil {
		newError("pack message").Base(err).WriteToLog()
		b.Release()
		return
	}
	b.Resize(0, int32(len(msgBytes)))

	if err := writer.WriteMessage(b); err != nil {
		newError("write IP answer").Base(err).WriteToLog()
	}
}

type outboundConn struct {
	access sync.Mutex
	dialer func() (internet.Connection, error)

	conn      net.Conn
	connReady chan struct{}
}

func (c *outboundConn) dial() error {
	conn, err := c.dialer()
	if err != nil {
		return err
	}
	c.conn = conn
	c.connReady <- struct{}{}
	return nil
}

func (c *outboundConn) Write(b []byte) (int, error) {
	c.access.Lock()

	if c.conn == nil {
		if err := c.dial(); err != nil {
			c.access.Unlock()
			newError("failed to dial outbound connection").Base(err).AtWarning().WriteToLog()
			return len(b), nil
		}
	}

	c.access.Unlock()

	return c.conn.Write(b)
}

func (c *outboundConn) Read(b []byte) (int, error) {
	var conn net.Conn
	c.access.Lock()
	conn = c.conn
	c.access.Unlock()

	if conn == nil {
		_, open := <-c.connReady
		if !open {
			return 0, io.EOF
		}
		conn = c.conn
	}

	return conn.Read(b)
}

func (c *outboundConn) Close() error {
	c.access.Lock()
	close(c.connReady)
	if c.conn != nil {
		c.conn.Close()
	}
	c.access.Unlock()
	return nil
}
