package assert

import (
	"strconv"
)

func Byte(value byte) *ByteSubject {
	return &ByteSubject{value: value}
}

type ByteSubject struct {
	Subject
	value byte
}

func (subject *ByteSubject) Named(name string) *ByteSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *ByteSubject) Fail(verb string, other byte) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + strconv.Itoa(int(other)) + ">.")
}

func (subject *ByteSubject) DisplayString() string {
	return subject.Subject.DisplayString(strconv.Itoa(int(subject.value)))
}

func (subject *ByteSubject) Equals(expectation byte) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *ByteSubject) GreaterThan(expectation byte) {
	if subject.value <= expectation {
		subject.Fail("is greater than", expectation)
	}
}

func (subject *ByteSubject) LessThan(expectation byte) {
	if subject.value >= expectation {
		subject.Fail("is less than", expectation)
	}
}
