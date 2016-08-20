package assert

import (
	"v2ray.com/core/common/serial"
)

func (this *Assert) Int64(value int64) *Int64Subject {
	return &Int64Subject{
		Subject: Subject{
			a:    this,
			disp: serial.Int64ToString(value),
		},
		value: value,
	}
}

type Int64Subject struct {
	Subject
	value int64
}

func (subject *Int64Subject) Equals(expectation int64) {
	if subject.value != expectation {
		subject.Fail("is equal to", serial.Int64ToString(expectation))
	}
}

func (subject *Int64Subject) GreaterThan(expectation int64) {
	if subject.value <= expectation {
		subject.Fail("is greater than", serial.Int64ToString(expectation))
	}
}

func (subject *Int64Subject) AtMost(expectation int64) {
	if subject.value > expectation {
		subject.Fail("is at most", serial.Int64ToString(expectation))
	}
}

func (subject *Int64Subject) AtLeast(expectation int64) {
	if subject.value < expectation {
		subject.Fail("is at least", serial.Int64ToString(expectation))
	}
}
