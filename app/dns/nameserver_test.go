package dns_test

import (
	"context"
	"testing"
	"time"

	. "v2ray.com/core/app/dns"
	. "v2ray.com/ext/assert"
)

func TestLocalNameServer(t *testing.T) {
	assert := With(t)

	s := NewLocalNameServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	ips, err := s.QueryIP(ctx, "google.com")
	cancel()
	assert(err, IsNil)
	assert(len(ips), GreaterThan, 0)
}
