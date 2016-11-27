package assert

import (
	"bytes"
	"net"
)

func (v *Assert) IP(value net.IP) *IPSubject {
	return &IPSubject{
		Subject: Subject{
			a:    v,
			disp: value.String(),
		},
		value: value,
	}
}

type IPSubject struct {
	Subject
	value net.IP
}

func (subject *IPSubject) IsNil() {
	if subject.value != nil {
		subject.Fail("is", "nil")
	}
}

func (subject *IPSubject) Equals(ip net.IP) {
	if !bytes.Equal([]byte(subject.value), []byte(ip)) {
		subject.Fail("equals to", ip.String())
	}
}
