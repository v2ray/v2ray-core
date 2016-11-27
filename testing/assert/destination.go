package assert

import (
	v2net "v2ray.com/core/common/net"
)

func (v *Assert) Destination(value v2net.Destination) *DestinationSubject {
	return &DestinationSubject{
		Subject: Subject{
			disp: value.String(),
			a:    v,
		},
		value: value,
	}
}

type DestinationSubject struct {
	Subject
	value v2net.Destination
}

func (v *DestinationSubject) IsTCP() {
	if v.value.Network != v2net.Network_TCP {
		v.Fail("is", "a TCP destination")
	}
}

func (v *DestinationSubject) IsNotTCP() {
	if v.value.Network == v2net.Network_TCP {
		v.Fail("is not", "a TCP destination")
	}
}

func (v *DestinationSubject) IsUDP() {
	if v.value.Network != v2net.Network_UDP {
		v.Fail("is", "a UDP destination")
	}
}

func (v *DestinationSubject) IsNotUDP() {
	if v.value.Network == v2net.Network_UDP {
		v.Fail("is not", "a UDP destination")
	}
}

func (v *DestinationSubject) EqualsString(another string) {
	if v.value.String() != another {
		v.Fail("not equals to string", another)
	}
}

func (v *DestinationSubject) HasAddress() *AddressSubject {
	return v.a.Address(v.value.Address)
}

func (v *DestinationSubject) HasPort() *PortSubject {
	return v.a.Port(v.value.Port)
}
