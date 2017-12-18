package http

import (
	"net/http"
	"strings"

	"v2ray.com/core/common/net"
)

func ParseXForwardedFor(header http.Header) []net.Address {
	xff := header.Get("X-Forwarded-For")
	if len(xff) == 0 {
		return nil
	}
	list := strings.Split(xff, ",")
	addrs := make([]net.Address, 0, len(list))
	for _, proxy := range list {
		addrs = append(addrs, net.ParseAddress(proxy))
	}
	return addrs
}
