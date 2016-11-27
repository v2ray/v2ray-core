package assert

import (
	"v2ray.com/core/common/serial"
)

func (v *Assert) Uint16(value uint16) *Uint16Subject {
	return &Uint16Subject{
		Subject: Subject{
			a:    v,
			disp: serial.Uint16ToString(value),
		},
		value: value,
	}
}

type Uint16Subject struct {
	Subject
	value uint16
}

func (subject *Uint16Subject) Equals(expectation uint16) {
	if subject.value != expectation {
		subject.Fail("is equal to", serial.Uint16ToString(expectation))
	}
}

func (subject *Uint16Subject) GreaterThan(expectation uint16) {
	if subject.value <= expectation {
		subject.Fail("is greater than", serial.Uint16ToString(expectation))
	}
}

func (subject *Uint16Subject) LessThan(expectation uint16) {
	if subject.value >= expectation {
		subject.Fail("is less than", serial.Uint16ToString(expectation))
	}
}

func (subject *Uint16Subject) IsPositive() {
	if subject.value <= 0 {
		subject.Fail("is", "positive")
	}
}

func (subject *Uint16Subject) IsNegative() {
	if subject.value >= 0 {
		subject.Fail("is not", "negative")
	}
}
