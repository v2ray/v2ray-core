package assert

import (
	"bytes"

	"fmt"

	"v2ray.com/core/common/serial"
)

func (v *Assert) Bytes(value []byte) *BytesSubject {
	return &BytesSubject{
		Subject: Subject{
			disp: serial.BytesToHexString(value),
			a:    v,
		},
		value: value,
	}
}

type BytesSubject struct {
	Subject
	value []byte
}

func (subject *BytesSubject) Equals(expectation []byte) {
	if len(subject.value) != len(expectation) {
		subject.FailWithMessage(fmt.Sprint("Bytes arrays have differen size: expected ", len(expectation), ", actual ", len(subject.value)))
		return
	}
	for idx, b := range expectation {
		if subject.value[idx] != b {
			subject.FailWithMessage(fmt.Sprint("Bytes are different: ", b, " vs ", subject.value[idx], " at pos ", idx))
			return
		}
	}
}

func (subject *BytesSubject) NotEquals(expectation []byte) {
	if bytes.Equal(subject.value, expectation) {
		subject.Fail("is not equal to", serial.BytesToHexString(expectation))
	}
}
