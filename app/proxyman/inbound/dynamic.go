package inbound

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
)

type DynamicInboundHandler struct {
	sync.Mutex
	tag            string
	ctx            context.Context
	cancel         context.CancelFunc
	proxyConfig    interface{}
	receiverConfig *proxyman.ReceiverConfig
	portsInUse     map[v2net.Port]bool
	worker         []worker
	worker2Recycle []worker
	lastRefresh    time.Time
}

func NewDynamicInboundHandler(ctx context.Context, tag string, receiverConfig *proxyman.ReceiverConfig, proxyConfig interface{}) (*DynamicInboundHandler, error) {
	ctx, cancel := context.WithCancel(ctx)
	h := &DynamicInboundHandler{
		ctx:            ctx,
		tag:            tag,
		cancel:         cancel,
		proxyConfig:    proxyConfig,
		receiverConfig: receiverConfig,
		portsInUse:     make(map[v2net.Port]bool),
	}

	return h, nil
}

func (h *DynamicInboundHandler) allocatePort() v2net.Port {
	from := int(h.receiverConfig.PortRange.From)
	delta := int(h.receiverConfig.PortRange.To) - from + 1
	h.Lock()
	defer h.Unlock()

	for {
		r := dice.Roll(delta)
		port := v2net.Port(from + r)
		_, used := h.portsInUse[port]
		if !used {
			h.portsInUse[port] = true
			return port
		}
	}
}

func (h *DynamicInboundHandler) refresh() error {
	h.lastRefresh = time.Now()

	ports2Del := make([]v2net.Port, 0, 16)
	for _, worker := range h.worker2Recycle {
		worker.Close()
		ports2Del = append(ports2Del, worker.Port())
	}

	h.Lock()
	for _, port := range ports2Del {
		delete(h.portsInUse, port)
	}
	h.Unlock()

	h.worker2Recycle, h.worker = h.worker, h.worker2Recycle[:0]

	for i := uint32(0); i < h.receiverConfig.AllocationStrategy.GetConcurrencyValue(); i++ {
		port := h.allocatePort()
		p, err := proxy.CreateInboundHandler(h.ctx, h.proxyConfig)
		if err != nil {
			log.Warning("Proxyman|DefaultInboundHandler: Failed to create proxy instance: ", err)
			continue
		}
		nl := p.Network()
		if nl.HasNetwork(v2net.Network_TCP) {
			worker := &tcpWorker{
				tag:              h.tag,
				address:          h.receiverConfig.Listen.AsAddress(),
				port:             port,
				proxy:            p,
				stream:           h.receiverConfig.StreamSettings,
				recvOrigDest:     h.receiverConfig.ReceiveOriginalDestination,
				allowPassiveConn: h.receiverConfig.AllowPassiveConnection,
			}
			if err := worker.Start(); err != nil {
				return err
			}
			h.worker = append(h.worker, worker)
		}

		if nl.HasNetwork(v2net.Network_UDP) {
			worker := &udpWorker{
				tag:          h.tag,
				proxy:        p,
				address:      h.receiverConfig.Listen.AsAddress(),
				port:         port,
				recvOrigDest: h.receiverConfig.ReceiveOriginalDestination,
			}
			if err := worker.Start(); err != nil {
				return err
			}
			h.worker = append(h.worker, worker)
		}
	}

	return nil
}

func (h *DynamicInboundHandler) monitor() {
	for {
		select {
		case <-h.ctx.Done():
			return
		case <-time.After(time.Minute * time.Duration(h.receiverConfig.AllocationStrategy.GetRefreshValue())):
			h.refresh()
		}
	}
}

func (h *DynamicInboundHandler) Start() error {
	err := h.refresh()
	go h.monitor()
	return err
}

func (h *DynamicInboundHandler) Close() {
	h.cancel()
}

func (h *DynamicInboundHandler) GetRandomInboundProxy() (proxy.InboundHandler, v2net.Port, int) {
	w := h.worker[dice.Roll(len(h.worker))]
	expire := h.receiverConfig.AllocationStrategy.GetRefreshValue() - uint32(time.Since(h.lastRefresh)/time.Minute)
	return w.Proxy(), w.Port(), int(expire)
}
