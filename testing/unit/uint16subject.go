package unit

import (
	"strconv"
)

type Uint16Subject struct {
	*Subject
	value uint16
}

func NewUint16Subject(base *Subject, value uint16) *Uint16Subject {
	subject := new(Uint16Subject)
	subject.Subject = base
	subject.value = value
	return subject
}

func (subject *Uint16Subject) Named(name string) *Uint16Subject {
	subject.Subject.Named(name)
	return subject
}

func (subject *Uint16Subject) Fail(verb string, other uint16) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + strconv.Itoa(int(other)) + ">.")
}

func (subject *Uint16Subject) DisplayString() string {
	return subject.Subject.DisplayString(strconv.Itoa(int(subject.value)))
}

func (subject *Uint16Subject) Equals(expectation uint16) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *Uint16Subject) GreaterThan(expectation uint16) {
	if subject.value <= expectation {
		subject.Fail("is greater than", expectation)
	}
}

func (subject *Uint16Subject) LessThan(expectation uint16) {
	if subject.value >= expectation {
		subject.Fail("is less than", expectation)
	}
}

func (subject *Uint16Subject) Positive() {
	if subject.value <= 0 {
		subject.FailWithMessage("Not true that " + subject.DisplayString() + " is positive.")
	}
}
