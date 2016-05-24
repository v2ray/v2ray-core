package assert

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

func (this *Assert) Port(value v2net.Port) *PortSubject {
	return &PortSubject{
		Subject: Subject{
			a:    this,
			disp: value.String(),
		},
		value: value,
	}
}

type PortSubject struct {
	Subject
	value v2net.Port
}

func (subject *PortSubject) Equals(expectation v2net.Port) {
	if subject.value.Value() != expectation.Value() {
		subject.Fail("is equal to", expectation.String())
	}
}

func (subject *PortSubject) GreaterThan(expectation v2net.Port) {
	if subject.value.Value() <= expectation.Value() {
		subject.Fail("is greater than", expectation.String())
	}
}

func (subject *PortSubject) LessThan(expectation v2net.Port) {
	if subject.value.Value() >= expectation.Value() {
		subject.Fail("is less than", expectation.String())
	}
}

func (subject *PortSubject) IsValid() {
	if subject.value == 0 {
		subject.Fail("is", "a valid port")
	}
}
