package assert

import (
	"strconv"
)

func Int64(value int64) *Int64Subject {
	return &Int64Subject{value: value}
}

type Int64Subject struct {
	Subject
	value int64
}

func (subject *Int64Subject) Named(name string) *Int64Subject {
	subject.Subject.Named(name)
	return subject
}

func (subject *Int64Subject) Fail(verb string, other int64) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + strconv.FormatInt(other, 10) + ">.")
}

func (subject *Int64Subject) DisplayString() string {
	return subject.Subject.DisplayString(strconv.FormatInt(subject.value, 10))
}

func (subject *Int64Subject) Equals(expectation int64) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *Int64Subject) GreaterThan(expectation int64) {
	if subject.value <= expectation {
		subject.Fail("is greater than", expectation)
	}
}

func (subject *Int64Subject) AtMost(expectation int64) {
	if subject.value > expectation {
		subject.Fail("is at most", expectation)
	}
}

func (subject *Int64Subject) AtLeast(expectation int64) {
	if subject.value < expectation {
		subject.Fail("is at least", expectation)
	}
}
