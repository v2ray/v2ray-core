package assert

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

func (this *Assert) Address(value v2net.Address) *AddressSubject {
	return &AddressSubject{
		Subject: Subject{
			disp: value.String(),
			a:    this,
		},
		value: value,
	}
}

type AddressSubject struct {
	Subject
	value v2net.Address
}

func (subject *AddressSubject) NotEquals(another v2net.Address) {
	if subject.value.Equals(another) {
		subject.Fail("not equals to", another.String())
	}
}

func (subject *AddressSubject) Equals(another v2net.Address) {
	if !subject.value.Equals(another) {
		subject.Fail("equals to", another.String())
	}
}

func (subject *AddressSubject) NotEqualsString(another string) {
	if subject.value.String() == another {
		subject.Fail("not equals to string", another)
	}
}

func (subject *AddressSubject) EqualsString(another string) {
	if subject.value.String() != another {
		subject.Fail("equals to string", another)
	}
}

func (subject *AddressSubject) IsIPv4() {
	if !subject.value.Family().IsIPv4() {
		subject.Fail("is", "an IPv4 address")
	}
}

func (subject *AddressSubject) IsNotIPv4() {
	if subject.value.Family().IsIPv4() {
		subject.Fail("is not", "an IPv4 address")
	}
}

func (subject *AddressSubject) IsIPv6() {
	if !subject.value.Family().IsIPv6() {
		subject.Fail("is", "an IPv6 address")
	}
}

func (subject *AddressSubject) IsNotIPv6() {
	if subject.value.Family().IsIPv6() {
		subject.Fail("is not", "an IPv6 address")
	}
}

func (subject *AddressSubject) IsDomain() {
	if !subject.value.Family().IsDomain() {
		subject.Fail("is", "a domain address")
	}
}

func (subject *AddressSubject) IsNotDomain() {
	if subject.value.Family().IsDomain() {
		subject.Fail("is not", "a domain address")
	}
}
