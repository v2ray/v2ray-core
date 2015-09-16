package unit

import (
	"strconv"
)

type BoolSubject struct {
	*Subject
	value bool
}

func NewBoolSubject(base *Subject, value bool) *BoolSubject {
	return &BoolSubject{
		Subject: base,
		value:   value,
	}
}

func (subject *BoolSubject) Named(name string) *BoolSubject {
	subject.Subject.Named(name)
	return subject
}

func (subject *BoolSubject) Fail(verb string, other bool) {
	subject.FailWithMessage("Not true that " + subject.DisplayString() + " " + verb + " <" + strconv.FormatBool(other) + ">.")
}

func (subject *BoolSubject) DisplayString() string {
	return subject.Subject.DisplayString(strconv.FormatBool(subject.value))
}

func (subject *BoolSubject) Equals(expectation bool) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

func (subject *BoolSubject) IsTrue() {
	if subject.value != true {
		subject.Fail("is", true)
	}
}

func (subject *BoolSubject) IsFalse() {
	if subject.value != false {
		subject.Fail("is", false)
	}
}
