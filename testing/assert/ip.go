package assert

import (
	"bytes"
	"net"

	"github.com/v2ray/v2ray-core/common/serial"
)

func IP(value net.IP) *IPSubject {
	return &IPSubject{value: value}
}

type IPSubject struct {
	Subject
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
		subject.Fail(subject.DisplayString(), "is", serial.StringT("nil"))
	}
}

func (subject *IPSubject) Equals(ip net.IP) {
	if !bytes.Equal([]byte(subject.value), []byte(ip)) {
		subject.Fail(subject.DisplayString(), "equals to", ip)
	}
}
