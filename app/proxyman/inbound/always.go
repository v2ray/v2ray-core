package inbound

import (
	"context"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type AlwaysOnInboundHandler struct {
	proxy   proxy.Inbound
	workers []worker
	mux     *mux.Server
	tag     string
}

func NewAlwaysOnInboundHandler(ctx context.Context, tag string, receiverConfig *proxyman.ReceiverConfig, proxyConfig interface{}) (*AlwaysOnInboundHandler, error) {
	rawProxy, err := common.CreateObject(ctx, proxyConfig)
	if err != nil {
		return nil, err
	}
	p, ok := rawProxy.(proxy.Inbound)
	if !ok {
		return nil, newError("not an inbound proxy.")
	}

	h := &AlwaysOnInboundHandler{
		proxy: p,
		mux:   mux.NewServer(ctx),
		tag:   tag,
	}

	nl := p.Network()
	pr := receiverConfig.PortRange
	address := receiverConfig.Listen.AsAddress()
	if address == nil {
		address = net.AnyIP
	}
	for port := pr.From; port <= pr.To; port++ {
		if nl.HasNetwork(net.Network_TCP) {
			newError("creating stream worker on ", address, ":", port).AtDebug().WriteToLog()
			worker := &tcpWorker{
				address:      address,
				port:         net.Port(port),
				proxy:        p,
				stream:       receiverConfig.StreamSettings,
				recvOrigDest: receiverConfig.ReceiveOriginalDestination,
				tag:          tag,
				dispatcher:   h.mux,
				sniffers:     receiverConfig.DomainOverride,
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
				dispatcher:   h.mux,
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

func (h *AlwaysOnInboundHandler) Close() error {
	for _, worker := range h.workers {
		worker.Close()
	}
	h.mux.Close()
	return nil
}

func (h *AlwaysOnInboundHandler) GetRandomInboundProxy() (interface{}, net.Port, int) {
	if len(h.workers) == 0 {
		return nil, 0, 0
	}
	w := h.workers[dice.Roll(len(h.workers))]
	return w.Proxy(), w.Port(), 9999
}

func (h *AlwaysOnInboundHandler) Tag() string {
	return h.tag
}

func (h *AlwaysOnInboundHandler) GetInbound() proxy.Inbound {
	return h.proxy
}
