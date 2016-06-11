package transport

import "github.com/v2ray/v2ray-core/transport/hub/kcpv"

var (
	connectionReuse = true
	enableKcp       = false
	KcpConfig       *kcpv.Config
)

func IsConnectionReusable() bool {
	return connectionReuse
}

func IsKcpEnabled() bool {
	return enableKcp
}
