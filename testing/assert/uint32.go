package assert

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

func (this *Assert) Uint32(value uint32) *Uint32Subject {
	return &Uint32Subject{
		Subject: Subject{
			a:    this,
			disp: serial.Uint32ToString(value),
		},
		value: value,
	}
}

type Uint32Subject struct {
	Subject
	value uint32
}

func (subject *Uint32Subject) Equals(expectation uint32) {
	if subject.value != expectation {
		subject.Fail("is equal to", serial.Uint32ToString(expectation))
	}
}

func (subject *Uint32Subject) GreaterThan(expectation uint32) {
	if subject.value <= expectation {
		subject.Fail("is greater than", serial.Uint32ToString(expectation))
	}
}

func (subject *Uint32Subject) LessThan(expectation uint32) {
	if subject.value >= expectation {
		subject.Fail("is less than", serial.Uint32ToString(expectation))
	}
}

func (subject *Uint32Subject) IsPositive() {
	if subject.value <= 0 {
		subject.Fail("is", "positive")
	}
}

func (subject *Uint32Subject) IsNegative() {
	if subject.value >= 0 {
		subject.Fail("is not", "negative")
	}
}
