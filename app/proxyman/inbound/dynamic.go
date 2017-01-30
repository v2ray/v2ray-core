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
	tag            string
	ctx            context.Context
	cancel         context.CancelFunc
	proxyConfig    interface{}
	receiverConfig *proxyman.ReceiverConfig
	portMutex      sync.Mutex
	portsInUse     map[v2net.Port]bool
	workerMutex    sync.RWMutex
	worker         []worker
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

	h.portMutex.Lock()
	defer h.portMutex.Unlock()

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

func (h *DynamicInboundHandler) waitAnyCloseWorkers(ctx context.Context, cancel context.CancelFunc, workers []worker, duration time.Duration) {
	time.Sleep(duration)
	cancel()
	ports2Del := make([]v2net.Port, len(workers))
	for idx, worker := range workers {
		ports2Del[idx] = worker.Port()
		worker.Close()
	}

	h.portMutex.Lock()
	for _, port := range ports2Del {
		delete(h.portsInUse, port)
	}
	h.portMutex.Unlock()
}

func (h *DynamicInboundHandler) refresh() error {
	h.lastRefresh = time.Now()

	timeout := time.Minute * time.Duration(h.receiverConfig.AllocationStrategy.GetRefreshValue()) * 2
	concurrency := h.receiverConfig.AllocationStrategy.GetConcurrencyValue()
	ctx, cancel := context.WithTimeout(h.ctx, timeout)
	workers := make([]worker, 0, concurrency)

	address := h.receiverConfig.Listen.AsAddress()
	if address == nil {
		address = v2net.AnyIP
	}
	for i := uint32(0); i < concurrency; i++ {
		port := h.allocatePort()
		p, err := proxy.CreateInboundHandler(ctx, h.proxyConfig)
		if err != nil {
			log.Warning("Proxyman|DefaultInboundHandler: Failed to create proxy instance: ", err)
			continue
		}
		nl := p.Network()
		if nl.HasNetwork(v2net.Network_TCP) {
			worker := &tcpWorker{
				tag:              h.tag,
				address:          address,
				port:             port,
				proxy:            p,
				stream:           h.receiverConfig.StreamSettings,
				recvOrigDest:     h.receiverConfig.ReceiveOriginalDestination,
				allowPassiveConn: h.receiverConfig.AllowPassiveConnection,
			}
			if err := worker.Start(); err != nil {
				log.Warning("Proxyman:InboundHandler: Failed to create TCP worker: ", err)
				continue
			}
			workers = append(workers, worker)
		}

		if nl.HasNetwork(v2net.Network_UDP) {
			worker := &udpWorker{
				tag:          h.tag,
				proxy:        p,
				address:      address,
				port:         port,
				recvOrigDest: h.receiverConfig.ReceiveOriginalDestination,
			}
			if err := worker.Start(); err != nil {
				log.Warning("Proxyman:InboundHandler: Failed to create UDP worker: ", err)
				continue
			}
			workers = append(workers, worker)
		}
	}

	h.workerMutex.Lock()
	h.worker = workers
	h.workerMutex.Unlock()

	go h.waitAnyCloseWorkers(ctx, cancel, workers, timeout)

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

func (h *DynamicInboundHandler) GetRandomInboundProxy() (proxy.Inbound, v2net.Port, int) {
	h.workerMutex.RLock()
	defer h.workerMutex.RUnlock()

	w := h.worker[dice.Roll(len(h.worker))]
	expire := h.receiverConfig.AllocationStrategy.GetRefreshValue() - uint32(time.Since(h.lastRefresh)/time.Minute)
	return w.Proxy(), w.Port(), int(expire)
}
