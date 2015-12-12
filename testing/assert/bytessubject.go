package assert

import (
	"bytes"
	"fmt"
)

func Bytes(value []byte) *BytesSubject {
	return &BytesSubject{value: value}
}

type BytesSubject struct {
	Subject
	value []byte
}

func (subject *BytesSubject) Named(name string) *BytesSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *BytesSubject) Fail(verb string, other []byte) {
	otherString := fmt.Sprintf("%v", other)
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + otherString + ">.")
}

func (subject *BytesSubject) DisplayString() string {
	return subject.Subject.DisplayString(fmt.Sprintf("%v", subject.value))
}

func (subject *BytesSubject) Equals(expectation []byte) {
	if !bytes.Equal(subject.value, expectation) {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *BytesSubject) NotEquals(expectation []byte) {
	if bytes.Equal(subject.value, expectation) {
		subject.Fail("is not equal to", expectation)
	}
}
