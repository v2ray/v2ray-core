package dns

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg dns -path App,DNS

import (
	"net"

	"v2ray.com/core/app"
)

// A Server is a DNS server for responding DNS queries.
type Server interface {
	Get(domain string) []net.IP
}

// FromSpace fetches a DNS server from context.
func FromSpace(space app.Space) Server {
	app := space.GetApplication((*Server)(nil))
	if app == nil {
		return nil
	}
	return app.(Server)
}
