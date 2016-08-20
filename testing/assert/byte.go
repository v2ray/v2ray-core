package assert

import (
	"v2ray.com/core/common/serial"
)

func (this *Assert) Byte(value byte) *ByteSubject {
	return &ByteSubject{
		Subject: Subject{
			disp: serial.ByteToHexString(value),
			a:    this,
		},
		value: value,
	}
}

type ByteSubject struct {
	Subject
	value byte
}

func (subject *ByteSubject) Equals(expectation byte) {
	if subject.value != expectation {
		subject.Fail("is equal to", serial.ByteToHexString(expectation))
	}
}

func (subject *ByteSubject) GreaterThan(expectation byte) {
	if subject.value <= expectation {
		subject.Fail("is greater than", serial.ByteToHexString(expectation))
	}
}

func (subject *ByteSubject) LessThan(expectation byte) {
	if subject.value >= expectation {
		subject.Fail("is less than", serial.ByteToHexString(expectation))
	}
}
