package log

import (
	"sync/atomic"
	"unsafe"

	"v2ray.com/core/common/serial"
)

type Message interface {
	String() string
}

type Handler interface {
	Handle(msg Message)
}

type noOpHandler byte

func (noOpHandler) Handle(msg Message) {}

type GeneralMessage struct {
	Severity Severity
	Content  interface{}
}

func (m *GeneralMessage) String() string {
	return serial.Concat("[", m.Severity, "]: ", m.Content)
}

func (s Severity) SevererThan(another Severity) bool {
	return s <= another
}

func Record(msg Message) {
	h := (*Handler)(atomic.LoadPointer(&logHandler))
	(*h).Handle(msg)
}

var (
	logHandler unsafe.Pointer
)

func RegisterHandler(handler Handler) {
	if handler == nil {
		panic("Log handler is nil")
	}
	atomic.StorePointer(&logHandler, unsafe.Pointer(&handler))
}

func init() {
	RegisterHandler(noOpHandler(0))
}
