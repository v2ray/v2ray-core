package inbound

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

func getStatCounter(v *core.Instance, tag string) (core.StatCounter, core.StatCounter) {
	var uplinkCounter core.StatCounter
	var downlinkCounter core.StatCounter

	policy := v.PolicyManager()
	stats := v.Stats()
	if len(tag) > 0 && policy.ForSystem().Stats.InboundUplink {
		name := "inbound>>>" + tag + ">>>traffic>>>uplink"
		c, _ := core.GetOrRegisterStatCounter(stats, name)
		if c != nil {
			uplinkCounter = c
		}
	}
	if len(tag) > 0 && policy.ForSystem().Stats.InboundDownlink {
		name := "inbound>>>" + tag + ">>>traffic>>>downlink"
		c, _ := core.GetOrRegisterStatCounter(stats, name)
		if c != nil {
			downlinkCounter = c
		}
	}

	return uplinkCounter, downlinkCounter
}

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

	uplinkCounter, downlinkCounter := getStatCounter(core.MustFromContext(ctx), tag)

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
				address:         address,
				port:            net.Port(port),
				proxy:           p,
				stream:          receiverConfig.StreamSettings,
				recvOrigDest:    receiverConfig.ReceiveOriginalDestination,
				tag:             tag,
				dispatcher:      h.mux,
				sniffers:        receiverConfig.DomainOverride,
				uplinkCounter:   uplinkCounter,
				downlinkCounter: downlinkCounter,
			}
			h.workers = append(h.workers, worker)
		}

		if nl.HasNetwork(net.Network_UDP) {
			worker := &udpWorker{
				tag:             tag,
				proxy:           p,
				address:         address,
				port:            net.Port(port),
				recvOrigDest:    receiverConfig.ReceiveOriginalDestination,
				dispatcher:      h.mux,
				uplinkCounter:   uplinkCounter,
				downlinkCounter: downlinkCounter,
			}
			h.workers = append(h.workers, worker)
		}
	}

	return h, nil
}

// Start implements common.Runnable.
func (h *AlwaysOnInboundHandler) Start() error {
	for _, worker := range h.workers {
		if err := worker.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Close implements common.Closable.
func (h *AlwaysOnInboundHandler) Close() error {
	var errors []interface{}
	for _, worker := range h.workers {
		if err := worker.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	if err := h.mux.Close(); err != nil {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return newError("failed to close all resources").Base(newError(serial.Concat(errors...)))
	}
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
