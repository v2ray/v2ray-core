package assert

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func Destination(value v2net.Destination) *DestinationSubject {
	return &DestinationSubject{value: value}
}

type DestinationSubject struct {
	*assert.Subject
	value v2net.Destination
}

func (this *DestinationSubject) Named(name string) *DestinationSubject {
	this.Subject.Named(name)
	return this
}

func (this *DestinationSubject) DisplayString() string {
	return this.Subject.DisplayString(this.value.String())
}

func (this *DestinationSubject) IsTCP() {
	if !this.value.IsTCP() {
		this.Fail(this.DisplayString(), "is", serial.StringLiteral("a TCP destination"))
	}
}

func (this *DestinationSubject) IsNotTCP() {
	if this.value.IsTCP() {
		this.Fail(this.DisplayString(), "is not", serial.StringLiteral("a TCP destination"))
	}
}

func (this *DestinationSubject) IsUDP() {
	if !this.value.IsUDP() {
		this.Fail(this.DisplayString(), "is", serial.StringLiteral("a UDP destination"))
	}
}

func (this *DestinationSubject) IsNotUDP() {
	if this.value.IsUDP() {
		this.Fail(this.DisplayString(), "is not", serial.StringLiteral("a UDP destination"))
	}
}
