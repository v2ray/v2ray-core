package assert

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

func (this *Assert) Pointer(value interface{}) *PointerSubject {
	return &PointerSubject{
		Subject: Subject{
			a:    this,
			disp: serial.PointerToString(value),
		},
		value: value,
	}
}

type PointerSubject struct {
	Subject
	value interface{}
}

func (subject *PointerSubject) Equals(expectation interface{}) {
	if subject.value != expectation {
		subject.Fail("is equal to", serial.PointerToString(expectation))
	}
}

func (subject *PointerSubject) IsNil() {
	if subject.value != nil {
		subject.Fail("is", "nil")
	}
}

func (subject *PointerSubject) IsNotNil() {
	if subject.value == nil {
		subject.Fail("is not", "nil")
	}
}
