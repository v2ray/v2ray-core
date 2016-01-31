package internal

import (
	"net"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/collect"
	"github.com/v2ray/v2ray-core/common/serial"
)

type entry struct {
	domain     string
	ip         net.IP
	validUntil time.Time
}

func newEntry(domain string, ip net.IP) *entry {
	this := &entry{
		domain: domain,
		ip:     ip,
	}
	this.Extend()
	return this
}

func (this *entry) IsValid() bool {
	return this.validUntil.After(time.Now())
}

func (this *entry) Extend() {
	this.validUntil = time.Now().Add(time.Hour)
}

type DnsCache struct {
	cache  *collect.ValidityMap
	config *CacheConfig
}

func NewCache(config *CacheConfig) *DnsCache {
	cache := &DnsCache{
		cache:  collect.NewValidityMap(3600),
		config: config,
	}
	return cache
}

func (this *DnsCache) Add(context app.Context, domain string, ip net.IP) {
	callerTag := context.CallerTag()
	if !this.config.IsTrustedSource(serial.StringLiteral(callerTag)) {
		return
	}

	this.cache.Set(serial.StringLiteral(domain), newEntry(domain, ip))
}

func (this *DnsCache) Get(context app.Context, domain string) net.IP {
	if value := this.cache.Get(serial.StringLiteral(domain)); value != nil {
		return value.(*entry).ip
	}
	return nil
}
