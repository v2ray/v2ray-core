package server

import (
	"time"

	"v2ray.com/core/common/net"
)

type IPResult struct {
	IP  []net.IP
	TTL time.Duration
}

type Querier interface {
	QueryDomain(domain string) <-chan *IPResult
}

type UDPQuerier struct {
	server net.Destination
}
