package assert

import (
	"v2ray.com/core/common/net"
)

func (v *Assert) Address(value net.Address) *AddressSubject {
	return &AddressSubject{
		Subject: Subject{
			disp: value.String(),
			a:    v,
		},
		value: value,
	}
}

type AddressSubject struct {
	Subject
	value net.Address
}

func (subject *AddressSubject) NotEquals(another net.Address) {
	if subject.value == another {
		subject.Fail("not equals to", another.String())
	}
}

func (subject *AddressSubject) Equals(another net.Address) {
	if subject.value != another {
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
