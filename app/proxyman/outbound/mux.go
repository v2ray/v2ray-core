package outbound

import "v2ray.com/core/proxy"

type mux struct {
	proxy  proxy.Outbound
	dialer proxy.Dialer
}
