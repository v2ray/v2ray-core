package assert

import (
	"v2ray.com/core/common/serial"
)

func (this *Assert) Int(value int) *IntSubject {
	return &IntSubject{
		Subject: Subject{
			a:    this,
			disp: serial.IntToString(value),
		},
		value: value,
	}
}

type IntSubject struct {
	Subject
	value int
}

func (subject *IntSubject) Equals(expectation int) {
	if subject.value != expectation {
		subject.Fail("is equal to", serial.IntToString(expectation))
	}
}

func (subject *IntSubject) GreaterThan(expectation int) {
	if subject.value <= expectation {
		subject.Fail("is greater than", serial.IntToString(expectation))
	}
}

func (subject *IntSubject) LessThan(expectation int) {
	if subject.value >= expectation {
		subject.Fail("is less than", serial.IntToString(expectation))
	}
}
