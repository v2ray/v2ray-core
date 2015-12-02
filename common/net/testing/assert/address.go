package assert

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func Address(value v2net.Address) *AddressSubject {
	return &AddressSubject{value: value}
}

type AddressSubject struct {
	assert.Subject
	value v2net.Address
}

func (subject *AddressSubject) Named(name string) *AddressSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *AddressSubject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value.String())
}

func (subject *AddressSubject) IsIPv4() {
	if !subject.value.IsIPv4() {
		subject.Fail(subject.DisplayString(), "is", serial.StringLiteral("an IPv4 address"))
	}
}

func (subject *AddressSubject) IsNotIPv4() {
	if subject.value.IsIPv4() {
		subject.Fail(subject.DisplayString(), "is not", serial.StringLiteral("an IPv4 address"))
	}
}

func (subject *AddressSubject) IsIPv6() {
	if !subject.value.IsIPv6() {
		subject.Fail(subject.DisplayString(), "is", serial.StringLiteral("an IPv6 address"))
	}
}

func (subject *AddressSubject) IsNotIPv6() {
	if subject.value.IsIPv6() {
		subject.Fail(subject.DisplayString(), "is not", serial.StringLiteral("an IPv6 address"))
	}
}

func (subject *AddressSubject) IsDomain() {
	if !subject.value.IsDomain() {
		subject.Fail(subject.DisplayString(), "is", serial.StringLiteral("a domain address"))
	}
}

func (subject *AddressSubject) IsNotDomain() {
	if subject.value.IsDomain() {
		subject.Fail(subject.DisplayString(), "is not", serial.StringLiteral("a domain address"))
	}
}
