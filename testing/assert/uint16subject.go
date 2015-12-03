package assert

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

func Uint16(value uint16) *Uint16Subject {
	return &Uint16Subject{value: serial.Uint16Literal(value)}
}

type Uint16Subject struct {
	Subject
	value serial.Uint16Literal
}

func (subject *Uint16Subject) Named(name string) *Uint16Subject {
	subject.Subject.Named(name)
	return subject
}

func (subject *Uint16Subject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value.String())
}

func (subject *Uint16Subject) Equals(expectation uint16) {
	if subject.value.Value() != expectation {
		subject.Fail(subject.DisplayString(), "is equal to", serial.Uint16Literal(expectation))
	}
}

func (subject *Uint16Subject) GreaterThan(expectation uint16) {
	if subject.value.Value() <= expectation {
		subject.Fail(subject.DisplayString(), "is greater than", serial.Uint16Literal(expectation))
	}
}

func (subject *Uint16Subject) LessThan(expectation uint16) {
	if subject.value.Value() >= expectation {
		subject.Fail(subject.DisplayString(), "is less than", serial.Uint16Literal(expectation))
	}
}

func (subject *Uint16Subject) Positive() {
	if subject.value.Value() <= 0 {
		subject.FailWithMessage("Not true that " + subject.DisplayString() + " is positive.")
	}
}
