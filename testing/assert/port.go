package assert

import (
	"v2ray.com/core/common/net"
)

func (v *Assert) Port(value net.Port) *PortSubject {
	return &PortSubject{
		Subject: Subject{
			a:    v,
			disp: value.String(),
		},
		value: value,
	}
}

type PortSubject struct {
	Subject
	value net.Port
}

func (subject *PortSubject) Equals(expectation net.Port) {
	if subject.value.Value() != expectation.Value() {
		subject.Fail("is equal to", expectation.String())
	}
}

func (subject *PortSubject) GreaterThan(expectation net.Port) {
	if subject.value.Value() <= expectation.Value() {
		subject.Fail("is greater than", expectation.String())
	}
}

func (subject *PortSubject) LessThan(expectation net.Port) {
	if subject.value.Value() >= expectation.Value() {
		subject.Fail("is less than", expectation.String())
	}
}

func (subject *PortSubject) IsValid() {
	if subject.value == 0 {
		subject.Fail("is", "a valid port")
	}
}
