package log

import (
	"sync/atomic"
	"unsafe"

	"v2ray.com/core/common/serial"
)

// Message is the interface for all log messages.
type Message interface {
	String() string
}

// Handler is the interface for log handler.
type Handler interface {
	Handle(msg Message)
}

type noOpHandler byte

func (noOpHandler) Handle(msg Message) {}

// GeneralMessage is a general log message that can contain all kind of content.
type GeneralMessage struct {
	Severity Severity
	Content  interface{}
}

// String implements Message.
func (m *GeneralMessage) String() string {
	return serial.Concat("[", m.Severity, "]: ", m.Content)
}

func (s Severity) SevererThan(another Severity) bool {
	return s <= another
}

// Record writes a message into log stream.
func Record(msg Message) {
	h := (*Handler)(atomic.LoadPointer(&logHandler))
	(*h).Handle(msg)
}

var (
	logHandler unsafe.Pointer
)

// RegisterHandler register a new handler as current log handler. Previous registered handler will be discarded.
func RegisterHandler(handler Handler) {
	if handler == nil {
		panic("Log handler is nil")
	}
	atomic.StorePointer(&logHandler, unsafe.Pointer(&handler))
}

func init() {
	RegisterHandler(noOpHandler(0))
}
