package assert

import (
	"github.com/v2ray/v2ray-core/common/serial"
)

func Int(value int) *IntSubject {
	return &IntSubject{value: serial.IntLiteral(value)}
}

type IntSubject struct {
	Subject
	value serial.IntLiteral
}

func (subject *IntSubject) Named(name string) *IntSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *IntSubject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value.String())
}

func (subject *IntSubject) Equals(expectation int) {
	if subject.value.Value() != expectation {
		subject.Fail(subject.DisplayString(), "is equal to", serial.IntLiteral(expectation))
	}
}

func (subject *IntSubject) GreaterThan(expectation int) {
	if subject.value.Value() <= expectation {
		subject.Fail(subject.DisplayString(), "is greater than", serial.IntLiteral(expectation))
	}
}

func (subject *IntSubject) LessThan(expectation int) {
	if subject.value.Value() >= expectation {
		subject.Fail(subject.DisplayString(), "is less than", serial.IntLiteral(expectation))
	}
}
