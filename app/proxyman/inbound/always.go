package inbound

import (
	"context"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type AlwaysOnInboundHandler struct {
	proxy   proxy.Inbound
	workers []worker
}

func NewAlwaysOnInboundHandler(ctx context.Context, tag string, receiverConfig *proxyman.ReceiverConfig, proxyConfig interface{}) (*AlwaysOnInboundHandler, error) {
	p, err := proxy.CreateInboundHandler(ctx, proxyConfig)
	if err != nil {
		return nil, err
	}

	h := &AlwaysOnInboundHandler{
		proxy: p,
	}

	nl := p.Network()
	pr := receiverConfig.PortRange
	address := receiverConfig.Listen.AsAddress()
	if address == nil {
		address = net.AnyIP
	}
	for port := pr.From; port <= pr.To; port++ {
		if nl.HasNetwork(net.Network_TCP) {
			log.Debug("Proxyman|DefaultInboundHandler: creating tcp worker on ", address, ":", port)
			worker := &tcpWorker{
				address:          address,
				port:             net.Port(port),
				proxy:            p,
				stream:           receiverConfig.StreamSettings,
				recvOrigDest:     receiverConfig.ReceiveOriginalDestination,
				tag:              tag,
				allowPassiveConn: receiverConfig.AllowPassiveConnection,
			}
			h.workers = append(h.workers, worker)
		}

		if nl.HasNetwork(net.Network_UDP) {
			worker := &udpWorker{
				tag:          tag,
				proxy:        p,
				address:      address,
				port:         net.Port(port),
				recvOrigDest: receiverConfig.ReceiveOriginalDestination,
			}
			h.workers = append(h.workers, worker)
		}
	}

	return h, nil
}

func (h *AlwaysOnInboundHandler) Start() error {
	for _, worker := range h.workers {
		if err := worker.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (h *AlwaysOnInboundHandler) Close() {
	for _, worker := range h.workers {
		worker.Close()
	}
}

func (h *AlwaysOnInboundHandler) GetRandomInboundProxy() (proxy.Inbound, net.Port, int) {
	w := h.workers[dice.Roll(len(h.workers))]
	return w.Proxy(), w.Port(), 9999
}
