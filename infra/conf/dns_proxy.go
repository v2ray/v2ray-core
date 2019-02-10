package conf

import (
	"github.com/golang/protobuf/proto"
	"v2ray.com/core/proxy/dns"
)

type DnsOutboundConfig struct{}

func (c *DnsOutboundConfig) Build() (proto.Message, error) {
	return new(dns.Config), nil
}
