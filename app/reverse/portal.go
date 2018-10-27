package reverse

import (
	"context"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/mux"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/task"
	"v2ray.com/core/common/vio"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/transport/pipe"
)

type Portal struct {
	ohm    outbound.Manager
	tag    string
	domain string
	picker *StaticMuxPicker
	client *mux.ClientManager
}

func NewPortal(config *PortalConfig, ohm outbound.Manager) (*Portal, error) {
	if len(config.Tag) == 0 {
		return nil, newError("portal tag is empty")
	}

	if len(config.Domain) == 0 {
		return nil, newError("portal domain is empty")
	}

	picker, err := NewStaticMuxPicker()
	if err != nil {
		return nil, err
	}

	return &Portal{
		ohm:    ohm,
		tag:    config.Tag,
		domain: config.Domain,
		picker: picker,
		client: &mux.ClientManager{
			Picker: picker,
		},
	}, nil
}

func (p *Portal) Start() error {
	return p.ohm.AddHandler(context.Background(), &Outbound{
		portal: p,
	})
}

func (p *Portal) Close() error {
	return p.ohm.RemoveHandler(context.Background(), p.tag)
}

func (s *Portal) HandleConnection(ctx context.Context, link *vio.Link) error {
	outboundMeta := session.OutboundFromContext(ctx)
	if outboundMeta == nil {
		return newError("outbound metadata not found").AtError()
	}

	if isInternalDomain(outboundMeta.Target) {
		muxClient, err := mux.NewClientWorker(*link, mux.ClientStrategy{
			MaxConcurrency: 0,
			MaxConnection:  256,
		})
		if err != nil {
			return newError("failed to create mux client worker").Base(err).AtWarning()
		}

		worker, err := NewPortalWorker(muxClient)
		if err != nil {
			return newError("failed to create portal worker").Base(err)
		}

		s.picker.AddWorker(worker)
		return nil
	}

	return s.client.Dispatch(ctx, link)
}

type Outbound struct {
	portal *Portal
	tag    string
}

func (o *Outbound) Tag() string {
	return o.tag
}

func (o *Outbound) Dispatch(ctx context.Context, link *vio.Link) {
	if err := o.portal.HandleConnection(ctx, link); err != nil {
		newError("failed to process reverse connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
		pipe.CloseError(link.Writer)
	}
}

func (o *Outbound) Start() error {
	return nil
}

func (o *Outbound) Close() error {
	return nil
}

type StaticMuxPicker struct {
	access  sync.Mutex
	workers []*PortalWorker
	cTask   *task.Periodic
}

func NewStaticMuxPicker() (*StaticMuxPicker, error) {
	p := &StaticMuxPicker{}
	p.cTask = &task.Periodic{
		Execute:  p.cleanup,
		Interval: time.Second * 30,
	}
	p.cTask.Start()
	return p, nil
}

func (p *StaticMuxPicker) cleanup() error {
	p.access.Lock()
	defer p.access.Unlock()

	var activeWorkers []*PortalWorker
	for _, w := range p.workers {
		if !w.Closed() {
			activeWorkers = append(activeWorkers, w)
		}
	}

	if len(activeWorkers) != len(p.workers) {
		p.workers = activeWorkers
	}

	return nil
}

func (p *StaticMuxPicker) PickAvailable() (*mux.ClientWorker, error) {
	p.access.Lock()
	defer p.access.Unlock()

	n := len(p.workers)
	if n == 0 {
		return nil, newError("empty worker list")
	}

	idx := dice.Roll(n)
	for i := 0; i < n; i++ {
		w := p.workers[(i+idx)%n]
		if !w.IsFull() {
			return w.client, nil
		}
	}

	return nil, newError("no mux client worker available")
}

func (p *StaticMuxPicker) AddWorker(worker *PortalWorker) {
	p.access.Lock()
	defer p.access.Unlock()

	p.workers = append(p.workers, worker)
}

type PortalWorker struct {
	client  *mux.ClientWorker
	control *task.Periodic
	writer  buf.Writer
	reader  buf.Reader
}

func NewPortalWorker(client *mux.ClientWorker) (*PortalWorker, error) {
	opt := []pipe.Option{pipe.WithSizeLimit(16 * 1024)}
	uplinkReader, uplinkWriter := pipe.New(opt...)
	downlinkReader, downlinkWriter := pipe.New(opt...)

	f := client.Dispatch(context.Background(), &vio.Link{
		Reader: uplinkReader,
		Writer: downlinkWriter,
	})
	if !f {
		return nil, newError("unable to dispatch control connection")
	}
	w := &PortalWorker{
		client: client,
		reader: downlinkReader,
		writer: uplinkWriter,
	}
	w.control = &task.Periodic{
		Execute:  w.heartbeat,
		Interval: time.Second * 2,
	}
	w.control.Start()
	return w, nil
}

func (w *PortalWorker) heartbeat() error {
	if w.client.Closed() {
		return newError("client worker stopped")
	}

	if w.writer == nil {
		return newError("already disposed")
	}

	msg := &Control{}
	msg.FillInRandom()

	if w.client.IsClosing() {
		msg.State = Control_DRAIN

		defer func() {
			common.Close(w.writer)
			pipe.CloseError(w.reader)
			w.writer = nil
		}()
	}

	b, err := proto.Marshal(msg)
	common.Must(err)
	var mb buf.MultiBuffer
	common.Must2(mb.Write(b))
	return w.writer.WriteMultiBuffer(mb)
}

func (w *PortalWorker) IsFull() bool {
	return w.client.IsFull()
}

func (w *PortalWorker) Closed() bool {
	return w.client.Closed()
}
