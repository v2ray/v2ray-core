package assert

import (
	v2net "v2ray.com/core/common/net"
)

func (this *Assert) Destination(value v2net.Destination) *DestinationSubject {
	return &DestinationSubject{
		Subject: Subject{
			disp: value.String(),
			a:    this,
		},
		value: value,
	}
}

type DestinationSubject struct {
	Subject
	value v2net.Destination
}

func (this *DestinationSubject) IsTCP() {
	if this.value.Network() != v2net.Network_TCP {
		this.Fail("is", "a TCP destination")
	}
}

func (this *DestinationSubject) IsNotTCP() {
	if this.value.Network() == v2net.Network_TCP {
		this.Fail("is not", "a TCP destination")
	}
}

func (this *DestinationSubject) IsUDP() {
	if this.value.Network() != v2net.Network_UDP {
		this.Fail("is", "a UDP destination")
	}
}

func (this *DestinationSubject) IsNotUDP() {
	if this.value.Network() == v2net.Network_UDP {
		this.Fail("is not", "a UDP destination")
	}
}

func (this *DestinationSubject) EqualsString(another string) {
	if this.value.String() != another {
		this.Fail("not equals to string", another)
	}
}

func (this *DestinationSubject) HasAddress() *AddressSubject {
	return this.a.Address(this.value.Address())
}

func (this *DestinationSubject) HasPort() *PortSubject {
	return this.a.Port(this.value.Port())
}
