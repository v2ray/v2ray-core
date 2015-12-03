package assert

import (
	"strconv"
)

// Assert on a boolean variable.
func Bool(value bool) *BoolSubject {
	return &BoolSubject{value: value}
}

type BoolSubject struct {
	Subject
	value bool
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

// to be equal to another boolean variable.
func (subject *BoolSubject) Equals(expectation bool) {
	if subject.value != expectation {
		subject.Fail("is equal to", expectation)
	}
}

// to be true.
func (subject *BoolSubject) IsTrue() {
	if subject.value != true {
		subject.Fail("is", true)
	}
}

// to be false.
func (subject *BoolSubject) IsFalse() {
	if subject.value != false {
		subject.Fail("is", false)
	}
}
