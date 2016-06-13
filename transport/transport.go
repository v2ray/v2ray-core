package transport

import "github.com/v2ray/v2ray-core/transport/hub/kcpv"

var (
	connectionReuse = true
	enableKcp       = false
	KcpConfig       *kcpv.Config
)

// IsConnectionReusable returns true if V2Ray is trying to reuse TCP connections.
func IsConnectionReusable() bool {
	return connectionReuse
}

func IsKcpEnabled() bool {
	return enableKcp
}
