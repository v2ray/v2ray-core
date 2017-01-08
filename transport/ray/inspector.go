package ray

import (
	"sync"

	"v2ray.com/core/common/buf"
)

type Inspector interface {
	Input(*buf.Buffer)
}

type NoOpInspector struct{}

func (NoOpInspector) Input(*buf.Buffer) {}

type InspectorChain struct {
	sync.RWMutex
	chain []Inspector
}

func (ic *InspectorChain) AddInspector(inspector Inspector) {
	ic.Lock()
	defer ic.Unlock()

	ic.chain = append(ic.chain, inspector)
}

func (ic *InspectorChain) Input(b *buf.Buffer) {
	ic.RLock()
	defer ic.RUnlock()

	for _, inspector := range ic.chain {
		inspector.Input(b)
	}
}
