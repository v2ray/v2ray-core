package http_test

import (
	"net/http"
	"testing"

	. "v2ray.com/core/common/protocol/http"
	. "v2ray.com/ext/assert"
)

func TestParseXForwardedFor(t *testing.T) {
	assert := With(t)

	header := http.Header{}
	header.Add("X-Forwarded-For", "129.78.138.66, 129.78.64.103")
	addrs := ParseXForwardedFor(header)
	assert(len(addrs), Equals, 2)
	assert(addrs[0].String(), Equals, "129.78.138.66")
	assert(addrs[1].String(), Equals, "129.78.64.103")
}
