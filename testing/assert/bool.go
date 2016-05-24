package assert

import (
	"strconv"
)

// Assert on a boolean variable.
func (this *Assert) Bool(value bool) *BoolSubject {
	return &BoolSubject{
		Subject: Subject{
			disp: strconv.FormatBool(value),
			a:    this,
		},
		value: value,
	}
}

type BoolSubject struct {
	Subject
	value bool
}

// to be equal to another boolean variable.
func (subject *BoolSubject) Equals(expectation bool) {
	if subject.value != expectation {
		subject.Fail("is equal to", strconv.FormatBool(expectation))
	}
}

// to be true.
func (subject *BoolSubject) IsTrue() {
	if subject.value != true {
		subject.Fail("is", "True")
	}
}

// to be false.
func (subject *BoolSubject) IsFalse() {
	if subject.value != false {
		subject.Fail("is", "False")
	}
}
