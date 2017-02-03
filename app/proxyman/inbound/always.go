package inbound

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type AlwaysOnInboundHandler struct {
	proxy      proxy.Inbound
	workers    []worker
	dispatcher dispatcher.Interface
}

func NewAlwaysOnInboundHandler(ctx context.Context, tag string, receiverConfig *proxyman.ReceiverConfig, proxyConfig interface{}) (*AlwaysOnInboundHandler, error) {
	p, err := proxy.CreateInboundHandler(ctx, proxyConfig)
	if err != nil {
		return nil, err
	}

	h := &AlwaysOnInboundHandler{
		proxy: p,
	}

	space := app.SpaceFromContext(ctx)
	space.OnInitialize(func() error {
		d := dispatcher.FromSpace(space)
		if d == nil {
			return errors.New("Proxyman|DefaultInboundHandler: No dispatcher in space.")
		}
		h.dispatcher = d
		return nil
	})

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
				dispatcher:       h,
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
				dispatcher:   h,
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

func (h *AlwaysOnInboundHandler) Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error) {
	return h.dispatcher.Dispatch(ctx, dest)
}
