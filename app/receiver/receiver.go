package receiver

import (
	"net"

	"v2ray.com/core/proxy"
)

type refresher struct {
}

type StreamReceiver struct {
	config *StreamReceiverConfig
	proxy  *proxy.InboundHandler

	listeners []net.Listener
}

func (r *StreamReceiver) Start() {

}
