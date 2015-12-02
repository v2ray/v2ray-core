package assert

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func Port(value v2net.Port) *PortSubject {
	return &PortSubject{value: value}
}

type PortSubject struct {
	assert.Subject
	value v2net.Port
}

func (subject *PortSubject) Named(name string) *PortSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *PortSubject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value.String())
}

func (subject *PortSubject) Equals(expectation v2net.Port) {
	if subject.value.Value() != expectation.Value() {
		subject.Fail(subject.DisplayString(), "is equal to", expectation)
	}
}

func (subject *PortSubject) GreaterThan(expectation v2net.Port) {
	if subject.value.Value() <= expectation.Value() {
		subject.Fail(subject.DisplayString(), "is greater than", expectation)
	}
}

func (subject *PortSubject) LessThan(expectation v2net.Port) {
	if subject.value.Value() >= expectation.Value() {
		subject.Fail(subject.DisplayString(), "is less than", expectation)
	}
}

func (subject *PortSubject) IsValid() {
  if subject.value == 0 {
    subject.Fail(subject.DisplayString(), "is", serial.StringLiteral("a valid port"))
  }
}
