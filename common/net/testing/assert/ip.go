package assert

import (
	"bytes"
	"net"

	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func IP(value net.IP) *IPSubject {
	return &IPSubject{value: value}
}

type IPSubject struct {
	assert.Subject
	value net.IP
}

func (subject *IPSubject) Named(name string) *IPSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *IPSubject) DisplayString() string {
	return subject.Subject.DisplayString(subject.value.String())
}

func (subject *IPSubject) IsNil() {
	if subject.value != nil {
		subject.Fail(subject.DisplayString(), "is", serial.StringLiteral("nil"))
	}
}

func (subject *IPSubject) Equals(ip net.IP) {
	if !bytes.Equal([]byte(subject.value), []byte(ip)) {
		subject.Fail(subject.DisplayString(), "equals to", ip)
	}
}
